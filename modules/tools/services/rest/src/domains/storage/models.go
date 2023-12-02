package storage

import (
	"time"

	bucketGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/storage/bucket"
	fileGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/storage/file"
)

type formatedBucket struct {
	Namespace string `json:"namespace"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`

	Hidden bool `json:"hidden"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Version   int64  `json:"version"`
}

func FormatedBuceketFromGRPC(bucket *bucketGRPC.Bucket) formatedBucket {
	return formatedBucket{
		Namespace: bucket.Namespace,
		UUID:      bucket.Uuid,
		Name:      bucket.Name,

		Hidden: bucket.Hidden,

		CreatedAt: bucket.XCreated.AsTime().Format(time.RFC3339),
		UpdatedAt: bucket.XUpdated.AsTime().Format(time.RFC3339),
		Version:   bucket.XVersion,
	}
}

type formatedFile struct {
	Namespace string `json:"namespace"`
	Bucket    string `json:"bucket"`
	UUID      string `json:"uuid"`
	Path      string `json:"path"`

	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Version   int64  `json:"version"`
}

func FormatedFileFromGRPC(file *fileGRPC.File) formatedFile {
	return formatedFile{
		Namespace: file.Namespace,
		Bucket:    file.Bucket,
		UUID:      file.Uuid,
		Path:      file.Path,

		MimeType: file.MimeType,
		Size:     file.Size,

		CreatedAt: file.XCreated.AsTime().Format(time.RFC3339),
		UpdatedAt: file.XUpdated.AsTime().Format(time.RFC3339),
		Version:   file.XVersion,
	}
}
