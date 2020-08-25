#!/bin/bash

echo "git pull"
git pull origin master

echo "build..."
go build

echo "重启服务"
pm2 restart zhuque

