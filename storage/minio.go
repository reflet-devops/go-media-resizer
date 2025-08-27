package storage

import (
	builtinCtx "context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jonboulle/clockwork"
	libMinio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/notification"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/mapstructure"
	"github.com/reflet-devops/go-media-resizer/types"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	MinioKey = "minio"
)

func init() {
	TypeStorageMapping[MinioKey] = createMinioStorage
}

var _ types.Storage = &minio{}

type ConfigMinio struct {
	ConfigClientMinio `mapstructure:",squash"`

	HealthCheckInterval time.Duration      `mapstructure:"health_check_interval" validate:"required"`
	FallbackMinio       *ConfigClientMinio `mapstructure:"fallback" validate:"required_unless=FallbackMinio nil"`
	PrefixPath          string             `mapstructure:"prefix_path"`
}

type ConfigClientMinio struct {
	Endpoint   string `mapstructure:"endpoint" validate:"required"`
	BucketName string `mapstructure:"bucket" validate:"required"`
	AccessKey  string `mapstructure:"access_key" validate:"required"`
	SecretKey  string `mapstructure:"secret_key" validate:"required"`
	UseSSL     bool   `mapstructure:"use_ssl"`
}

type minio struct {
	currentClient   types.MinioClient
	primaryClient   types.MinioClient
	secondaryClient types.MinioClient
	cfg             ConfigMinio

	ctx   *context.Context
	mx    sync.RWMutex
	clock clockwork.Clock
}

func (m *minio) Type() string {
	return MinioKey
}

func (m *minio) getFullPath(path string) string {
	path = strings.TrimLeft(path, "/")
	if m.cfg.PrefixPath == "" {
		return path
	}
	return strings.Join([]string{m.cfg.PrefixPath, path}, "/")
}

func (m *minio) GetFile(path string) (io.Reader, error) {
	object, err := m.getClient().GetObject(builtinCtx.Background(), m.cfg.BucketName, filepath.Join(m.cfg.PrefixPath, path), libMinio.GetObjectOptions{})

	if err != nil {
		return nil, err
	}

	stat, errStat := object.Stat()
	if errStat != nil {
		return nil, errStat
	}

	if stat.Size == 0 {
		return nil, os.ErrNotExist
	}

	return object, err
}

func (m *minio) getClient() types.MinioClient {
	m.mx.RLock()
	defer m.mx.RUnlock()
	return m.currentClient
}

func (m *minio) startFallback() {
	cancel, errHc := m.primaryClient.HealthCheck(time.Second * 2)
	if errHc != nil {
		defer cancel()
	}
	failedCount := 0
	ticker := m.clock.NewTicker(m.cfg.HealthCheckInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.Chan():
			if m.primaryClient.IsOnline() {
				failedCount = 0
				if m.currentClient != m.primaryClient {
					m.ctx.Logger.Info(fmt.Sprintf("use minio primary as currentClient %s/%s", m.cfg.Endpoint, m.cfg.BucketName))
					m.mx.Lock()
					m.currentClient = m.primaryClient
					m.mx.Unlock()
				}
			} else {
				failedCount++
				m.ctx.Logger.Error(fmt.Sprintf("primary minio is offline %s/%s", m.cfg.Endpoint, m.cfg.BucketName))
				if m.currentClient != m.secondaryClient && failedCount >= 3 {
					m.ctx.Logger.Error(fmt.Sprintf("use minio secondary as currentClient %s/%s", m.cfg.Endpoint, m.cfg.BucketName))
					m.mx.Lock()
					m.currentClient = m.secondaryClient
					m.mx.Unlock()
				}
			}
		case <-m.ctx.Done():
			return
		}
	}
}

func (m *minio) NotifyFileChange(chanEvent chan types.Events) {
	eventsTypes := []string{
		string(notification.ObjectCreatedPut),
		string(notification.ObjectCreatedPost),
		string(notification.ObjectCreatedCopy),
		string(notification.ObjectRemovedDelete),
	}
	minioChanEvent := m.primaryClient.ListenBucketNotification(
		builtinCtx.Background(),
		m.cfg.BucketName,
		m.getFullPath("*"),
		"",
		eventsTypes,
	)
	for {
		select {
		case minioEvent := <-minioChanEvent:
			if minioEvent.Err != nil {
				m.ctx.Logger.Error(fmt.Sprintf("minio notify failed: %v", minioEvent.Err))
				continue
			}
			events := types.Events{}
			for _, record := range minioEvent.Records {
				event := types.Event{
					Type: types.EventTypePurge,
					Path: strings.Replace(record.S3.Object.Key, m.getFullPath(""), "", 1),
				}
				events = append(events, event)
			}
			chanEvent <- events
		case <-m.ctx.Done():
			return
		}
	}
}

func createMinioStorage(ctx *context.Context, cfg config.StorageConfig) (types.Storage, error) {
	instanceConfig := ConfigMinio{
		HealthCheckInterval: time.Second * 5,
	}

	err := mapstructure.Decode(cfg.Config, &instanceConfig)
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	err = validate.Struct(instanceConfig)
	if err != nil {
		return nil, err
	}

	instanceConfig.PrefixPath = strings.Trim(instanceConfig.PrefixPath, "/")

	minioClient, errNewClient := createMinioClient(instanceConfig.ConfigClientMinio)
	if errNewClient != nil {
		return nil, errNewClient
	}

	instance := &minio{currentClient: minioClient, primaryClient: minioClient, cfg: instanceConfig, ctx: ctx, clock: clockwork.NewRealClock()}

	if instanceConfig.FallbackMinio != nil {
		minioClient2, errNewClient2 := createMinioClient(*instanceConfig.FallbackMinio)
		if errNewClient2 != nil {
			return nil, errNewClient2
		}
		instance.secondaryClient = minioClient2
		go instance.startFallback()
	}

	return instance, nil
}

func createMinioClient(cfg ConfigClientMinio) (*libMinio.Client, error) {
	return libMinio.New(
		cfg.Endpoint,
		&libMinio.Options{
			Creds:        credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure:       cfg.UseSSL,
			BucketLookup: libMinio.BucketLookupAuto,
		},
	)
}
