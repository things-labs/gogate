#!/bin/bash

echo "building..."
dataTime=`date +%F.%T.%z`
CGO_ENABLED=1 GOOS=linux GOARCH=mipsle GOMIPS=softfloat \
STAGING_DIR=/opt/toolchain/openwrt18.06/staging_dir \
CC=/opt/toolchain/openwrt18.06/staging_dir/gcc-mipsel-linux-7.3.0/bin/mipsel-openwrt-linux-gcc \
go build -ldflags "-X github.com/thinkgos/gogate/misc.BuildTime=$dataTime -s -w" -o gogate-mipsle .
if [ $? -ne 0 ]
then
	echo "go build failed"
	exit
fi

bzip2 -c gogate-mipsle > gogate-mipsle.bz2
if [ $? -eq 0 ]
then
    echo "build success"
else
    echo "build failed"
fi
