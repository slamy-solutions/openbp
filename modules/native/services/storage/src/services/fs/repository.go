package fs

import (
	"context"
	"crypto/rand"
	crypto "crypto/subtle"
	"errors"
	"io"
	"log/slog"
	"strings"
	"time"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FileRepository struct {
	systemStub *system.SystemStub
	logger     *slog.Logger
}

func NewFSRepository(systemStub *system.SystemStub, logger *slog.Logger) (*FileRepository, error) {
	err := prepareCollections(context.Background(), systemStub, "")
	if err != nil {
		return nil, errors.Join(errors.New("failed to prepare files collection"), err)
	}

	return &FileRepository{
		systemStub: systemStub,
		logger:     logger.With("repository", "file"),
	}, nil
}

func generateDownloadSecret(length int) (string, error) {
	chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes), nil
}

func (r *FileRepository) getBaseDirectoryPath(path string) string {
	pathSplit := strings.Split(strings.TrimLeft(path, "/"), "/")
	if len(pathSplit) == 1 {
		return "/"
	}

	return "/" + strings.Join(pathSplit[:len(pathSplit)-1], "/")
}

func (r *FileRepository) isPathValid(path string) bool {
	return !strings.Contains(path, "//") && strings.HasPrefix(path, "/")
}

func (r *FileRepository) Create(ctx context.Context, namespace string, bucket primitive.ObjectID, path string, mimeType string) (*File, error) {
	if !r.isPathValid(path) {
		return nil, ErrFilePathInvalid
	}

	err := r.MkDirs(ctx, namespace, bucket, r.getBaseDirectoryPath(path))
	if err != nil {
		err = errors.Join(errors.New("failed to create directories"), err)
		return nil, err
	}

	collection := GetFileInfoCollection(r.systemStub, namespace)
	creationTime := time.Now().UTC()

	downloadSecret, err := generateDownloadSecret(32)
	if err != nil {
		err = errors.Join(errors.New("failed to generate download secret"), err)
		r.logger.Error("Failed to generate download secret", "error", err.Error())
		return nil, err
	}

	file := File{
		Namespace:            namespace,
		Bucket:               bucket,
		Path:                 path,
		BaseDirectoryPath:    r.getBaseDirectoryPath(path),
		MimeType:             mimeType,
		Size:                 0,
		DirectDownloadSecret: downloadSecret,
		GridFSFile:           primitive.NilObjectID,
		Created:              creationTime,
		Updated:              creationTime,
		Version:              0,
	}
	result, err := collection.InsertOne(ctx, file)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrFileAlreadyExists
		}

		err = errors.Join(errors.New("failed to insert file info"), err)
		r.logger.Error("Failed to insert file info", "error", err.Error())
		return nil, err
	}

	file.UUID = result.InsertedID.(primitive.ObjectID)
	return &file, nil
}

func (r *FileRepository) Upload(ctx context.Context, namespace string, fileInfoUUID primitive.ObjectID, file io.Reader) (*File, error) {
	collection := GetFileInfoCollection(r.systemStub, namespace)
	var oldFile File
	err := collection.FindOne(ctx, bson.M{"_id": fileInfoUUID}).Decode(&oldFile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrFileNotFound
		}

		err = errors.Join(errors.New("failed to find file info"), err)
		r.logger.Error("Failed to find file info", "error", err.Error())
		return nil, err
	}

	gridFSBucket, err := GetFileBucket(r.systemStub, namespace, oldFile.Bucket)
	if err != nil {
		r.logger.Error("Failed to get file bucket", "error", err.Error())
		return nil, err
	}

	uploadStream, err := gridFSBucket.OpenUploadStreamWithID(fileInfoUUID, "file")
	if err != nil {
		err = errors.Join(errors.New("failed to open upload stream"), err)
		r.logger.Error("Failed to open upload stream", "error", err.Error())
		return nil, err
	}

	fileSize, err := io.Copy(uploadStream, file)
	if err != nil {
		uploadStream.Abort()

		err = errors.Join(errors.New("failed to copy file to upload stream"), err)
		r.logger.Error("Failed to copy file to upload stream", "error", err.Error())
		return nil, err
	}

	err = uploadStream.Close()
	if err != nil {
		uploadStream.Abort()

		err = errors.Join(errors.New("failed to close upload stream"), err)
		r.logger.Error("Failed to close upload stream", "error", err.Error())
		return nil, err
	}

	//Update file info
	var fileInfo File
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": fileInfoUUID},
		bson.M{
			"$set":         bson.M{"gridfsFile": uploadStream.FileID, "size": fileSize},
			"$inc":         bson.M{"_version": 1},
			"$currentDate": bson.M{"_updated": bson.M{"$type": "timestamp"}},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&fileInfo)
	if err != nil {
		deleteErr := gridFSBucket.Delete(uploadStream.FileID)
		if deleteErr != nil {
			r.logger.Warn("Failed to delete file from gridfs after problems with inserting file info", "error", deleteErr.Error())
		}

		if err == mongo.ErrNoDocuments {
			return nil, ErrFileNotFound
		}

		err = errors.Join(errors.New("failed to update file info"), err)
		r.logger.Error("Failed to update file info", "error", err.Error())
		return nil, err
	}

	if oldFile.GridFSFile != primitive.NilObjectID {
		deleteErr := gridFSBucket.Delete(oldFile.GridFSFile)
		if deleteErr != nil {
			r.logger.Warn("Failed to delete old file from gridfs after inserting new file info", "error", deleteErr.Error())
		}
	}

	fileInfo.Namespace = namespace
	return &fileInfo, nil
}

func (r *FileRepository) Stat(ctx context.Context, namespace string, uuid primitive.ObjectID) (*File, error) {
	collection := GetFileInfoCollection(r.systemStub, namespace)

	var fileInfo File
	err := collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&fileInfo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrFileNotFound
		}

		err = errors.Join(errors.New("failed to find file info"), err)
		r.logger.Error("Failed to find file info", "error", err.Error())
		return nil, err
	}

	fileInfo.Namespace = namespace
	return &fileInfo, nil
}

func (r *FileRepository) Update(ctx context.Context, namespace string, uuid primitive.ObjectID, path string, mimeType string) (*File, error) {
	if !r.isPathValid(path) {
		return nil, ErrFilePathInvalid
	}

	collection := GetFileInfoCollection(r.systemStub, namespace)

	var oldFile File
	err := collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&oldFile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrFileNotFound
		}

		err = errors.Join(errors.New("failed to find file info"), err)
		r.logger.Error("Failed to find file info", "error", err.Error())
		return nil, err
	}

	err = r.MkDirs(ctx, namespace, oldFile.Bucket, r.getBaseDirectoryPath(path))
	if err != nil {
		err = errors.Join(errors.New("failed to create directories"), err)
		return nil, err
	}

	var fileInfo File
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": uuid},
		bson.M{
			"$set": bson.M{
				"path":              path,
				"baseDirectoryPath": r.getBaseDirectoryPath(path),
				"mimeType":          mimeType,
			},
			"$inc":         bson.M{"_version": 1},
			"$currentDate": bson.M{"_updated": bson.M{"$type": "timestamp"}},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&fileInfo)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrFileNotFound
		}

		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrFileAlreadyExists
		}

		err = errors.Join(errors.New("failed to update file info"), err)
		r.logger.Error("Failed to update file info", "error", err.Error())
		return nil, err
	}

	fileInfo.Namespace = namespace
	return &fileInfo, nil
}

func (r *FileRepository) Delete(ctx context.Context, namespace string, uuid primitive.ObjectID) (*File, error) {
	collection := GetFileInfoCollection(r.systemStub, namespace)

	var fileInfo File
	err := collection.FindOneAndDelete(ctx, bson.M{"_id": uuid}).Decode(&fileInfo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrFileNotFound
		}

		err = errors.Join(errors.New("failed to delete file info"), err)
		r.logger.Error("Failed to delete file info", "error", err.Error())
		return nil, err
	}

	if fileInfo.GridFSFile != primitive.NilObjectID {
		gridFSBucket, err := GetFileBucket(r.systemStub, namespace, fileInfo.Bucket)
		if err != nil {
			r.logger.Error("Failed to get file bucket", "error", err.Error())
			return nil, err
		}

		err = gridFSBucket.Delete(fileInfo.GridFSFile)
		if err != nil {
			r.logger.Warn("Failed to delete file from gridfs after deleting file info", "error", err.Error())
		}
	}

	fileInfo.Namespace = namespace
	return &fileInfo, nil
}

type FilesListCursor struct {
	ctx        context.Context
	namespace  string
	repository *FileRepository
	cursor     *mongo.Cursor
}

func (c *FilesListCursor) Next() (*File, error) {
	hasNext := c.cursor.Next(c.ctx)
	if !hasNext {
		if err := c.cursor.Err(); err != nil {
			err = errors.Join(errors.New("failed to get next file from database"), err)
			c.repository.logger.Error("failed to list files", "error", err, slog.String("namespace", c.namespace))
			return nil, err
		}

		return nil, io.EOF
	}

	var file File
	err := c.cursor.Decode(file)
	if err != nil {
		err = errors.Join(errors.New("failed to decode file from the badatase"), err)
		c.repository.logger.Error("Failed to decode file", "error", err, slog.String("namespace", c.namespace))
		return nil, err
	}

	file.Namespace = c.namespace
	return &file, nil
}
func (c *FilesListCursor) Close() error {
	return c.cursor.Close(context.Background())
}

func (r *FileRepository) List(ctx context.Context, namespace string, bucket primitive.ObjectID, skip int64, limit int64) (*FilesListCursor, error) {
	collection := GetFileInfoCollection(r.systemStub, namespace)

	opts := options.Find()
	if skip > 0 {
		opts.SetSkip(skip)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}
	opts.SetSort(bson.M{"_id": 1})

	cursor, err := collection.Find(ctx, bson.M{"bucket": bucket}, opts)
	if err != nil {
		err = errors.Join(errors.New("failed to list files"), err)
		r.logger.Error("Failed to list files", "error", err.Error())
		return nil, err
	}

	return &FilesListCursor{
		ctx:        ctx,
		namespace:  namespace,
		repository: r,
		cursor:     cursor,
	}, nil
}

func (r *FileRepository) Count(ctx context.Context, namespace string, bucket primitive.ObjectID) (int64, error) {
	collection := GetFileInfoCollection(r.systemStub, namespace)

	count, err := collection.CountDocuments(ctx, bson.M{"bucket": bucket})
	if err != nil {
		err = errors.Join(errors.New("failed to count files"), err)
		r.logger.Error("Failed to count files", "error", err.Error())
		return 0, err
	}

	return count, nil
}

func (r *FileRepository) Download(ctx context.Context, namespace string, uuid primitive.ObjectID, seek int64) (*File, io.ReadCloser, error) {
	collection := GetFileInfoCollection(r.systemStub, namespace)

	var fileInfo File
	err := collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&fileInfo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, ErrFileNotFound
		}

		err = errors.Join(errors.New("failed to find file info"), err)
		r.logger.Error("Failed to find file info", "error", err.Error())
		return nil, nil, err
	}

	gridFSBucket, err := GetFileBucket(r.systemStub, namespace, fileInfo.Bucket)
	if err != nil {
		r.logger.Error("Failed to get file bucket", "error", err.Error())
		return nil, nil, err
	}

	downloadStream, err := gridFSBucket.OpenDownloadStream(fileInfo.GridFSFile)
	if err != nil {
		err = errors.Join(errors.New("failed to open download stream"), err)
		r.logger.Error("Failed to open download stream", "error", err.Error())
		return nil, nil, err
	}

	if seek != 0 {
		_, err = downloadStream.Skip(seek)
		if err != nil {
			err = errors.Join(errors.New("failed to seek download stream"), err)
			r.logger.Error("Failed to seek download stream", "error", err.Error())
			return nil, nil, err
		}
	}

	fileInfo.Namespace = namespace
	return &fileInfo, downloadStream, nil
}

func (r *FileRepository) DownloadByPath(ctx context.Context, namespace string, bucket primitive.ObjectID, path string, seek int64) (*File, io.ReadCloser, error) {
	collection := GetFileInfoCollection(r.systemStub, namespace)

	var fileInfo File
	err := collection.FindOne(ctx, bson.M{"bucket": bucket, "path": path}).Decode(&fileInfo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, ErrFileNotFound
		}

		err = errors.Join(errors.New("failed to find file info"), err)
		r.logger.Error("Failed to find file info", "error", err.Error())
		return nil, nil, err
	}

	gridFSBucket, err := GetFileBucket(r.systemStub, namespace, fileInfo.Bucket)
	if err != nil {
		r.logger.Error("Failed to get file bucket", "error", err.Error())
		return nil, nil, err
	}

	downloadStream, err := gridFSBucket.OpenDownloadStream(fileInfo.GridFSFile)
	if err != nil {
		err = errors.Join(errors.New("failed to open download stream"), err)
		r.logger.Error("Failed to open download stream", "error", err.Error())
		return nil, nil, err
	}

	if seek != 0 {
		_, err = downloadStream.Skip(seek)
		if err != nil {
			err = errors.Join(errors.New("failed to seek download stream"), err)
			r.logger.Error("Failed to seek download stream", "error", err.Error())
			return nil, nil, err
		}
	}

	fileInfo.Namespace = namespace
	return &fileInfo, downloadStream, nil
}

func (r *FileRepository) DownloadDirect(ctx context.Context, namespace string, uuid primitive.ObjectID, directSercret string, seek int64) (*File, io.ReadCloser, error) {
	collection := GetFileInfoCollection(r.systemStub, namespace)

	var fileInfo File
	err := collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&fileInfo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, ErrFileNotFound
		}

		err = errors.Join(errors.New("failed to find file info"), err)
		r.logger.Error("Failed to find file info", "error", err.Error())
		return nil, nil, err
	}

	if crypto.ConstantTimeCompare([]byte(fileInfo.DirectDownloadSecret), []byte(directSercret)) != 1 {
		return nil, nil, ErrDirectDownloadSecretInvalid
	}

	gridFSBucket, err := GetFileBucket(r.systemStub, namespace, fileInfo.Bucket)
	if err != nil {
		r.logger.Error("Failed to get file bucket", "error", err.Error())
		return nil, nil, err
	}

	downloadStream, err := gridFSBucket.OpenDownloadStream(fileInfo.GridFSFile)
	if err != nil {
		err = errors.Join(errors.New("failed to open download stream"), err)
		r.logger.Error("Failed to open download stream", "error", err.Error())
		return nil, nil, err
	}

	if seek != 0 {
		_, err = downloadStream.Skip(seek)
		if err != nil {
			err = errors.Join(errors.New("failed to seek download stream"), err)
			r.logger.Error("Failed to seek download stream", "error", err.Error())
			return nil, nil, err
		}
	}

	fileInfo.Namespace = namespace
	return &fileInfo, downloadStream, nil
}

func (r *FileRepository) DownloadDirectByPath(ctx context.Context, namespace string, bucket primitive.ObjectID, path string, directSercret string, seek int64) (*File, io.ReadCloser, error) {
	collection := GetFileInfoCollection(r.systemStub, namespace)

	var fileInfo File
	err := collection.FindOne(ctx, bson.M{"bucket": bucket, "path": path}).Decode(&fileInfo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, ErrFileNotFound
		}

		err = errors.Join(errors.New("failed to find file info"), err)
		r.logger.Error("Failed to find file info", "error", err.Error())
		return nil, nil, err
	}

	if crypto.ConstantTimeCompare([]byte(fileInfo.DirectDownloadSecret), []byte(directSercret)) != 1 {
		return nil, nil, ErrDirectDownloadSecretInvalid
	}

	gridFSBucket, err := GetFileBucket(r.systemStub, namespace, fileInfo.Bucket)
	if err != nil {
		r.logger.Error("Failed to get file bucket", "error", err.Error())
		return nil, nil, err
	}

	downloadStream, err := gridFSBucket.OpenDownloadStream(fileInfo.GridFSFile)
	if err != nil {
		err = errors.Join(errors.New("failed to open download stream"), err)
		r.logger.Error("Failed to open download stream", "error", err.Error())
		return nil, nil, err
	}

	if seek != 0 {
		_, err = downloadStream.Skip(seek)
		if err != nil {
			err = errors.Join(errors.New("failed to seek download stream"), err)
			r.logger.Error("Failed to seek download stream", "error", err.Error())
			return nil, nil, err
		}
	}

	fileInfo.Namespace = namespace
	return &fileInfo, downloadStream, nil
}

func (r *FileRepository) Ls(ctx context.Context, namespace string, bucket primitive.ObjectID, path string) (*FilesListCursor, error) {
	collection := GetFileInfoCollection(r.systemStub, namespace)

	opts := options.Find()
	opts.SetSort(bson.M{"_id": 1})

	cursor, err := collection.Find(ctx, bson.M{"bucket": bucket, "baseDirectoryPath": path}, opts)
	if err != nil {
		err = errors.Join(errors.New("failed to list files"), err)
		r.logger.Error("Failed to list files", "error", err.Error())
		return nil, err
	}

	return &FilesListCursor{
		ctx:        ctx,
		namespace:  namespace,
		repository: r,
		cursor:     cursor,
	}, nil
}

func (r *FileRepository) MkDirs(ctx context.Context, namespace string, bucket primitive.ObjectID, path string) error {
	if !r.isPathValid(path) {
		return ErrFilePathInvalid
	}

	if path == "/" {
		return nil
	}

	parentDirectoriesSplt := strings.Split(strings.TrimLeft(path, "/"), "/")
	parentDirectories := make([]string, 0, len(parentDirectoriesSplt))
	for i, _ := range parentDirectoriesSplt {
		parentDirectories = append(parentDirectories, "/"+strings.Join(parentDirectoriesSplt[:i+1], "/"))
	}

	collection := GetDirectoryCollection(r.systemStub, namespace)

	writeOperations := make([]mongo.WriteModel, 0, len(parentDirectories))
	for _, parentDirectory := range parentDirectories {
		directory := Directory{
			Namespace:         namespace,
			Bucket:            bucket,
			Path:              parentDirectory,
			BaseDirectoryPath: r.getBaseDirectoryPath(parentDirectory),
		}
		writeOperations = append(writeOperations, mongo.NewUpdateOneModel().SetFilter(bson.M{"bucket": bucket, "path": parentDirectory}).SetUpdate(bson.M{"$set": directory}).SetUpsert(true))
	}

	_, err := collection.BulkWrite(ctx, writeOperations)
	if err != nil {
		err = errors.Join(errors.New("failed to insert directory info"), err)
		r.logger.Error("Failed to insert directory info", "error", err.Error())
		return err
	}

	return nil
}
