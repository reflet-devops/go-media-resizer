#!/bin/bash

go install go.uber.org/mock/mockgen@latest

rm -rf mocks
mockgen -destination=mocks/afero/fs.go -package=mockAfero github.com/spf13/afero Fs
