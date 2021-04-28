#!/bin/bash

rm zhuque
echo "上传文件......"
rz
echo "文件执行权限"
chmod u+x zhuque
echo "git pull"
git pull origin master
echo "重启服务"
pm2 restart zhuque
pm2 logs
