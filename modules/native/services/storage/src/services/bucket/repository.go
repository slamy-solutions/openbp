package bucket

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"regexp"
	"time"

	"github.com/slamy-solutions/openbp/modules/native/services/storage/src/services/fs"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var bucketNameRegex *regexp.Regexp

func init() {
	var err error
	bucketNameRegex, err = regexp.Compile("^[a-zA-Z0-9-_]+$")
	if err != nil {
		panic(err)
	}
}

type BucketRepository struct {
	systemStub *system.SystemStub

	logger *slog.Logger
}

func NewBucketRepository(ctx context.Context, systemStub *system.SystemStub, logger *slog.Logger) (*BucketRepository, error) {
	err := prepareBucketsCollection(ctx, systemStub, "")
	if err != nil {
		return nil, errors.Join(errors.New("failed to prepare buckets collection"), err)
	}

	return &BucketRepository{
		systemStub: systemStub,
		logger:     logger.With("repository", "bucket"),
	}, nil
}

func (r *BucketRepository) Create(ctx context.Context, namespace string, name string, hidden bool) (*Bucket, error) {
	if !bucketNameRegex.MatchString(name) {
		return nil, ErrBucketNameInvalid
	}

	collection := GetBucketsCollection(r.systemStub, namespace)

	creationTime := time.Now().UTC()

	bucket := &Bucket{
		Namespace: namespace,
		Name:      name,
		Hidden:    hidden,
		Created:   creationTime,
		Updated:   creationTime,
		Version:   0,
	}

	result, err := collection.InsertOne(ctx, bucket)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrBucketAlreadyExists
		}

		err = errors.Join(errors.New("failed to insert new sbucket to the badatase"), err)
		r.logger.Error("Failed to create bucket", "error", err, bucket.ToSlogAttr("bucket"))
		return nil, err
	}
	bucket.UUID = result.InsertedID.(primitive.ObjectID)

	r.logger.Info("Bucket created", bucket.ToSlogAttr("bucket"))
	return bucket, nil
}

func (r *BucketRepository) Ensure(ctx context.Context, namespace string, name string, hidden bool) (*Bucket, error) {
	if !bucketNameRegex.MatchString(name) {
		return nil, ErrBucketNameInvalid
	}

	collection := GetBucketsCollection(r.systemStub, namespace)

	creationTime := time.Now().UTC()

	bucket := &Bucket{
		Namespace: namespace,
		Name:      name,
		Hidden:    hidden,
		Created:   creationTime,
		Updated:   creationTime,
		Version:   0,
	}

	err := collection.FindOneAndUpdate(
		ctx,
		bson.M{"name": name},
		bson.M{"$setOnInsert": bucket},
		options.FindOneAndUpdate().SetUpsert(true),
	).Decode(bucket)
	if err != nil {
		err = errors.Join(errors.New("failed to ensure bucket in the badatase"), err)
		r.logger.Error("Failed to ensure bucket", "error", err, bucket.ToSlogAttr("bucket"))
		return nil, err
	}

	if bucket.Created.After(creationTime.Add(time.Second * -1)) {
		r.logger.Info("Bucket created after ensuring operation", bucket.ToSlogAttr("bucket"))
	}

	bucket.Namespace = namespace
	return bucket, nil
}

func (r *BucketRepository) GetByUUID(ctx context.Context, namespace string, uuid primitive.ObjectID) (*Bucket, error) {
	collection := GetBucketsCollection(r.systemStub, namespace)

	bucket := &Bucket{}
	err := collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(bucket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrBucketNotFound
		}

		err = errors.Join(errors.New("failed to get bucket from the badatase"), err)
		r.logger.Error("Failed to get bucket", "error", err, slog.String("namespace", namespace), slog.String("uuid", uuid.Hex()))
		return nil, err
	}

	bucket.Namespace = namespace
	return bucket, nil
}

func (r *BucketRepository) Get(ctx context.Context, namespace string, name string) (*Bucket, error) {
	collection := GetBucketsCollection(r.systemStub, namespace)

	bucket := &Bucket{}
	err := collection.FindOne(ctx, bson.M{"name": name}).Decode(bucket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrBucketNotFound
		}

		err = errors.Join(errors.New("failed to get bucket from the badatase"), err)
		r.logger.Error("Failed to get bucket", "error", err, slog.String("namespace", namespace), slog.String("name", name))
		return nil, err
	}

	bucket.Namespace = namespace
	return bucket, nil
}

type BucketsListCursor struct {
	ctx        context.Context
	namespace  string
	repository *BucketRepository
	cursor     *mongo.Cursor
}

func (c *BucketsListCursor) Next() (*Bucket, error) {
	hasNext := c.cursor.Next(c.ctx)
	if !hasNext {
		if err := c.cursor.Err(); err != nil {
			err = errors.Join(errors.New("failed to get next bucket from database"), err)
			c.repository.logger.Error("failed to list buckets", "error", err, slog.String("namespace", c.namespace))
			return nil, err
		}

		return nil, io.EOF
	}

	var bucket Bucket
	err := c.cursor.Decode(bucket)
	if err != nil {
		err = errors.Join(errors.New("failed to decode bucket from the badatase"), err)
		c.repository.logger.Error("Failed to decode bucket", "error", err, slog.String("namespace", c.namespace))
		return nil, err
	}

	bucket.Namespace = c.namespace
	return &bucket, nil
}
func (c *BucketsListCursor) Close() error {
	return c.cursor.Close(context.Background())
}

func (r *BucketRepository) List(ctx context.Context, namespace string, skip int32, limit int32) (*BucketsListCursor, error) {
	collection := GetBucketsCollection(r.systemStub, namespace)

	opts := options.Find()
	if skip > 0 {
		opts.SetSkip(int64(skip))
	}
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	opts.SetSort(bson.M{"_id": 1})

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		err = errors.Join(errors.New("failed to list buckets from the badatase"), err)
		r.logger.Error("Failed to list buckets", "error", err, slog.String("namespace", namespace))
		return nil, err
	}

	return &BucketsListCursor{
		ctx:        ctx,
		namespace:  namespace,
		repository: r,
		cursor:     cursor,
	}, nil
}

func (r *BucketRepository) Count(ctx context.Context, namespace string) (int64, error) {
	collection := GetBucketsCollection(r.systemStub, namespace)
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		err = errors.Join(errors.New("failed to count buckets in the badatase"), err)
		r.logger.Error("Failed to count buckets", "error", err, slog.String("namespace", namespace))
		return 0, err
	}

	return count, nil
}

func (r *BucketRepository) Update(ctx context.Context, namespace string, uuid primitive.ObjectID, name string, hidden bool) (*Bucket, error) {
	if !bucketNameRegex.MatchString(name) {
		return nil, ErrBucketNameInvalid
	}

	collection := GetBucketsCollection(r.systemStub, namespace)

	var bucket Bucket
	err := collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": uuid},
		bson.M{
			"$set": bson.M{
				"name":   name,
				"hidden": hidden,
			},
			"$inc": bson.M{
				"_version": 1,
			},
			"$currentDate": bson.M{"_updated": bson.M{"$type": "timestamp"}},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&bucket)
	if err != nil {
		if err == mongo.ErrNilDocument {
			return nil, ErrBucketNotFound
		}

		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrBucketAlreadyExists
		}

		err = errors.Join(errors.New("failed to update bucket in the badatase"), err)
		r.logger.Error("Failed to update bucket", "error", err, slog.String("namespace", namespace))
		return nil, err
	}

	bucket.Namespace = namespace
	return &bucket, nil
}

func (r *BucketRepository) Delete(ctx context.Context, namespace string, name string) (*Bucket, error) {
	collection := GetBucketsCollection(r.systemStub, namespace)

	var bucket Bucket
	err := collection.FindOneAndDelete(ctx, bson.M{"name": name}).Decode(&bucket)
	if err != nil {
		if err == mongo.ErrNilDocument {
			return nil, ErrBucketNotFound
		}

		err = errors.Join(errors.New("failed to delete bucket from the badatase"), err)
		r.logger.Error("Failed to delete bucket", "error", err, slog.String("namespace", namespace))
		return nil, err
	}

	err = fs.DestroyCollectionsForBucket(ctx, r.systemStub, namespace, bucket.UUID)
	if err != nil {
		r.logger.Error("Failed to delete files data (bucket) while deleting bucket", "error", err, slog.String("namespace", namespace), slog.String("bucket", bucket.UUID.Hex()))
	}

	bucket.Namespace = namespace
	return &bucket, nil
}

func (r *BucketRepository) DeleteByUUID(ctx context.Context, namespace string, uuid primitive.ObjectID) (*Bucket, error) {
	collection := GetBucketsCollection(r.systemStub, namespace)

	var bucket Bucket
	err := collection.FindOneAndDelete(ctx, bson.M{"_id": uuid}).Decode(&bucket)
	if err != nil {
		if err == mongo.ErrNilDocument {
			return nil, ErrBucketNotFound
		}

		err = errors.Join(errors.New("failed to delete bucket from the badatase"), err)
		r.logger.Error("Failed to delete bucket", "error", err, slog.String("namespace", namespace))
		return nil, err
	}

	err = fs.DestroyCollectionsForBucket(ctx, r.systemStub, namespace, bucket.UUID)
	if err != nil {
		r.logger.Error("Failed to delete files data (bucket) while deleting bucket", "error", err, slog.String("namespace", namespace), slog.String("bucket", bucket.UUID.Hex()))
	}

	bucket.Namespace = namespace
	return &bucket, nil
}
