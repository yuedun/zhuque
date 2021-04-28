#!/bin/bash

rm zhuque
echo "上传文件......"
rz
echo "文件执行权限"
chmod u+x zhuque
echo "git pull origin $1"
git pull origin $1
echo "重启服务"
pm2 restart zhuque
pm2 logs
