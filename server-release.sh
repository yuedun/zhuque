#!/bin/bash

echo "git pull"
git pull origin master

echo "build..."
go build

export GIN_MODE=release
echo "重启服务"
pm2 restart zhuque

pm2 logs
