#!/bin/bash

if [ x"$1" = x ]; then
    echo "请输入分支参数"
    exit 1
fi

echo "git pull origin $1"
git pull origin $1

echo "build..."
go build

export GIN_MODE=release
echo "重启服务"
pm2 restart zhuque

pm2 logs
