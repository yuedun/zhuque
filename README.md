# 朱雀发布系统

<!-- TOC -->
- [朱雀发布系统](#朱雀发布系统)
    - [介绍](#介绍)
    - [使用框架文档](#使用框架文档)
    - [本地开发](#本地开发)
        - [依赖项](#依赖项)
        - [首次运行](#首次运行)
    - [部署到服务器](#部署到服务器)
        - [部署流程](#部署流程)
        - [系统要求](#系统要求)
        - [权限架构](#权限架构)
        - [简单部署](#简单部署)
        - [Complete tutorial](#complete-tutorial)
        - [生态系统文件参考](#生态系统文件参考)
        - [部署选项](#部署选项)
    - [SCP发布](#scp发布)
        - [发布流程](#发布流程)
        - [配置说明](#配置说明)
<!-- /TOC -->

#### [文档](https://yuedun.gitbook.io/zhuque/)
## 介绍
朱雀发布系统是专门为nodejs发布而开发的系统，使用go语言开发，前端使用了layui mini框架，简单易上手。
朱雀发布系统前后端分离，但又是一体的，无需部署两套服务。

理论上朱雀发布系统可以发布其他语言应用程序，但其与nodejs更加相得益彰。
原因是朱雀发布系统依赖于PM2，PM2大家都知道，与nodejs几乎是绝配，所以有nodejs的地方很大可能有PM2，所以使用朱雀发布系统就显得更加简单了，无需专门安装PM2。

（推荐）同时也支持scp（rsync）发布模式。优点是一键部署发布，使用简单。


数据库使用sqlite，无需单独安装和配置。

## 使用框架文档
[前端框架layuimini](http://layuimini.99php.cn/docs/index.html)

[前端框架layuimini](http://layuimini.99php.cn/onepage/v2/index.html)

[前端框架layui](https://www.layui.com/doc/)

[Font Awesome图标库](https://fontawesome.dashgame.com/)

[后端框架gorm](http://gorm.book.jasperxu.com/)
[后端框架gorm](https://gorm.io/zh_CN/docs/index.html)
[后端框架gorm](https://v1.gorm.io/zh_CN/docs/index.html)

[后端框架gin](https://github.com/gin-gonic/gin#using-middleware)

## 本地开发

### 依赖项
项目使用了sqlite3，需要安装gcc，参考地址：[gcc安装](https://www.jianshu.com/p/dc0fc5d8c900)
如果只想使用mysql的话注释掉代码`github.com/jinzhu/gorm/dialects/sqlite`即可。
### 首次运行
1. 复制`conf-sample.yaml`文件为`conf.yaml`文件。不做任何修改也可以运行。
2. 配置`env`变量为`debug`。可选
3. 初始化数据库，见`zhuque.sql`。
4. 使用`test`账号，密码`test`登录系统。


## 部署到服务器
### 部署流程
1. 在服务器指定位置下载源码。
2. [linux版本下载](https://github.com/yuedun/zhuque/releases/download/v1.3.0/zhuque)
或者安装go环境。go可以交叉编译，但是由于sqlite的缘故，windows环境下并不能顺利的编译linux版本，所以最好还是在linux环境下编译。
3. 配置`conf.yaml`文件，参照`conf-sample.yaml`文件。
4. 项目目录中编译`go build`，第一次会安装依赖会慢一些。
5. `./zhuque`启动服务。

### 系统要求
1核2G，2核4G都可以，运行时占用内存只有几M，所以对系统配置要求不高。
### 权限架构
该系统权限使用了基于角色的访问控制方法（RBAC）。是目前公认的解决大型企业的统一资源访问控制的有效方法。 其显著的两大特征是：1.减小授权管理的复杂性，降低管理开销。2.灵活地支持企业的安全策略，并对企业的变化有很大的伸缩性。

参考文档：

[官方原文](https://pm2.keymetrics.io/docs/usage/deployment/)

[PM2自动部署代码流程总结](https://segmentfault.com/a/1190000017310047)

[pm2 官方文档 学习笔记](https://my.oschina.net/u/4400196/blog/3283439)

[通过Github与PM2部署Node应用](https://zhuanlan.zhihu.com/p/20940096)

[deploy](https://github.com/Unitech/PM2/blob/0.14.7/ADVANCED_README.md#deployment-options)

### 简单部署

你只需要在ecosystem.json文件中添加**deploy**属性。 下面是是部署一个应用的最低要求:

process.json:

```
{
   "apps" : [{
      "name" : "HTTP-API",
      "script" : "http.js"
   }],
   "deploy" : {
     // "production" 是环境变量名
     "production" : {
       "user" : "ubuntu",
       "host" : ["192.168.0.13"],
       "ref"  : "origin/master",
       "repo" : "git@github.com:Username/repository.git",
       "path" : "/var/www/my-repository",
       "post-deploy" : "npm install; grunt dist"
      },
   }
}
```

/bin/bash:
```
# 部署到远程服务
$ pm2 deploy production setup

# 更新远程版本
$ pm2 deploy production update

# 回滚到上一版本
$ pm2 deploy production revert 1

# 在远程机器上执行命令
$ pm2 deploy production exec "pm2 reload all"
```

### Complete tutorial

1- 生成一个样本ecosystem.json列出进程和部署环境的文件。

`pm2 ecosystem`

在当前文件夹会创建一个`ecosystem.json`文件:
```
{
  // Applications part
  "apps" : [{
    "name"      : "API",
    "script"    : "app.js",
    "env": {
      "COMMON_VARIABLE": "true"
    },
    // Environment variables injected when starting with --env production
    // http://pm2.keymetrics.io/docs/usage/application-declaration/#switching-to-different-environments
    "env_production" : {
      "NODE_ENV": "production"
    }
  },{
    "name"      : "WEB",
    "script"    : "web.js"
  }],
  // 部署部分
  // Here you describe each environment
  "deploy" : {
    "production" : {
      "user" : "node",
      // 服务器集群
      "host" : ["212.83.163.1", "212.83.163.2", "212.83.163.3"],
      // 分支
      "ref"  : "origin/master",
      // Git 地址
      "repo" : "git@github.com:repo.git",
      // 应用在服务器上的地址
      "path" : "/var/www/production",
      // Can be used to give options in the format used in the configura-
      // tion file.  This is useful for specifying options for which there
      // is no separate command-line flag, see 'man ssh'
      // can be either a single string or an array of strings
      "ssh_options": "StrictHostKeyChecking=no",
      // To prepare the host by installing required software (eg: git)
      // even before the setup process starts
      // can be multiple commands separated by the character ";"
      // or path to a script on your local machine
      "pre-setup" : "apt-get install git",
      // Commands / path to a script on the host machine
      // This will be executed on the host after cloning the repository
      // eg: placing configurations in the shared dir etc
      "post-setup": "ls -la",
      // Commands to execute locally (on the same machine you deploy things)
      // Can be multiple commands separated by the character ";"
      "pre-deploy-local" : "echo 'This is a local executed command'"
      // Commands to be executed on the server after the repo has been cloned
      "post-deploy" : "npm install && pm2 startOrRestart ecosystem.json --env production"
      // Environment variables that must be injected in all applications on this env
      "env"  : {
        "NODE_ENV": "production"
      }
    },
    "staging" : {
      "user" : "node",
      "host" : "212.83.163.1",
      "ref"  : "origin/master",
      "repo" : "git@github.com:repo.git",
      "path" : "/var/www/development",
      "ssh_options": ["StrictHostKeyChecking=no", "PasswordAuthentication=no"],
      "post-deploy" : "pm2 startOrRestart ecosystem.json --env dev",
      "env"  : {
        "NODE_ENV": "staging"
      }
    }
  }
}
```

Edit the file according to your needs.

2- 确保您的本地机器上有ssh公钥
```
ssh-keygen -t rsa
ssh-copy-id node@myserver.com
```

If you encounter any errors, please refer to the troubleshooting section below.

3- 初始化远程项目:

`pm2 deploy <configuration_file> <environment> setup`

例如:

`pm2 deploy ecosystem.json production setup`

这个命令会在远程服务上创建文件夹。

4- 部署代码

`pm2 deploy ecosystem.json production`

Now your code will be populated, installed and started with PM2.

### 生态系统文件参考

生态系统文件的目的是收集应用所有的配置选项和环境变量。

### 部署选项

选项名称|描述|类型|默认
---|---|---|---
`key`|`SSH`密钥的路径|`String`|`$HOME/.ssh`
`user`|`SSH`用户|`String`|
`host`|`SSH`主机|`[String]`|
`ssh_options`|`SSH`选项，不包括命令行标志，查看`man ssh`|`String`，`[String]`|
`ref`|`GIT`的`remote/branch`|`String`|
`repo`|`GIT`的`remote`|`String`|
`path`|服务器中的路径|`String`|
`pre-setup`|1.远程机器拉代码之前|`String`|
`post-setup`|2.远程服务器拉代码|`String`|
`pre-deploy-local`|3.post-deploy前在宿主机上执行的命令|`String`|
`post-deploy`|4.部署后执行|`String`|

## SCP发布
scp发布是较为常用的发布方式，但朱雀使用的是rsync发布，主要是利用其增量同步功能，加速代码同步。

### 发布流程
推荐的做法是在发布机上拉代码，编译。
同步代码到远程应用服务器，重启。
也就是说远程应用服务只需重启即可，不需要做编译操作，在发布机上编译即可。

### 配置说明
```js
{
      "user": "root",
      "host": ["10.11.12.13"],
      "ref":"master",
      "repo": "git@github.com/yuedun/zhuque.git",
      "path": "/data/www/zhuque",
      "preBuild" : "",
      "build":"go build",
      "preDeploy" : "echo '发布前置';",
      "postDeploy" : "pm2 restart zheque;pm2 ls",
      "rsyncArgs":"-u --delete --exclude '.git' --exclude '.env'"
}
```


- user：用户名
- host[]：主机地址
- ref：分支
- repo：仓库地址
- path：项目部署路径，需要包含项目目录
- preBuild：编译前置，在发布机上编译代码设置的环境变量，比如前端项目设置不同的环境变量使用不同的接口地址。
- build：编译命令，例如：npm run build
- preDeploy：发布前置，应用服务重启前设置环境变量等操作。
- postDeploy：发布命令，例如：pm2 reload zhuque
- rsyncArgs：rsync参数

常用rsync参数：
- --exclude 排除不进行同步的文件，比如--exclude="*.iso"
- --delete参数删除只存在于目标目录、不存在于源目标的文件，即保证目标目录是源目标的镜像。
如果想要排除指定的文件，即不删除某个文件，可以使用exclude指定，例如：--exclude '.env'，会删除其他文件而不会删除.env文件。

