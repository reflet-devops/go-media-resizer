#!/bin/bash

export MC_HOST_minio="http://root:rootroot@127.0.0.1:9000"
export MC_INSECURE=true
export MC_DEBUG=false
export MC_JSON=true
mc ${@}
