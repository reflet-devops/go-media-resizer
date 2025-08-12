package types

import (
	"context"
	"github.com/minio/minio-go/v7"
	"time"
)

type MinioClient interface {
	//ListBuckets(ctx context.Context) ([]minio.BucketInfo, error)
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	HealthCheck(hcDuration time.Duration) (context.CancelFunc, error)
	IsOnline() bool
	//IsOffline() bool
}
