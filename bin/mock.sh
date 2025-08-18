#!/bin/bash

go install go.uber.org/mock/mockgen@latest

rm -rf mocks
mockgen -destination=mocks/afero/fs.go -package=mockAfero github.com/spf13/afero Fs
mockgen -destination=mocks/types/minio.go -package=mockTypes github.com/reflet-devops/go-media-resizer/types MinioClient
mockgen -destination=mocks/types/http.go -package=mockTypes github.com/reflet-devops/go-media-resizer/types Client

mockgen -destination=mocks/types/storage.go -package=mockTypes github.com/reflet-devops/go-media-resizer/types Storage,PurgeCache
