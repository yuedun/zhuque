#!/bin/bash

go build #本地开发使用
# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .. #发布到Linux环境下时使用该命令编译