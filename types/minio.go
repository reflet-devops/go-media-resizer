package types

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/notification"
	"time"
)

type MinioClient interface {
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	HealthCheck(hcDuration time.Duration) (context.CancelFunc, error)
	IsOnline() bool
	ListenBucketNotification(ctx context.Context, bucketName, prefix, suffix string, events []string) <-chan notification.Info
}
