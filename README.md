# 朱雀-nodejs发布系统
朱雀发布系统是专门为nodejs发布而开发的系统，使用go语言开发，前端使用了layui mini框架，简单易上手。
朱雀发布系统前后端分离，但又是一体的，无需部署两套服务。
理论上朱雀发布系统可以发布其他语言应用程序，但其与nodejs更加相得益彰。
原因是朱雀发布系统依赖于PM2，PM2大家都知道，与nodejs几乎是绝配，所以有nodejs的地方很大可能有PM2，所以使用朱雀发布系统就显得更加简单了，无需专门安装PM2。

（推荐）同时也支持scp（rsync）发布模式。优点是一键部署发布，使用简单。


数据库使用sqlite，无需单独安装和配置。

## 使用框架文档
[前端框架](http://layuimini.99php.cn/docs/index.html)
[前端框架](http://layuimini.99php.cn/onepage/v2/index.html)
[前端框架](https://www.layui.com/doc/)
[Font Awesome图标库](https://fontawesome.dashgame.com/)
[后端框架gorm](http://gorm.book.jasperxu.com/)
[后端框架gin](https://github.com/gin-gonic/gin#using-middleware)

## 本地开发

### 依赖项
项目使用了sqlite3，需要安装gcc，参考地址：[gcc安装](https://www.jianshu.com/p/dc0fc5d8c900)
### 首次运行
1. 复制`conf-sample.yaml`文件为`conf.yaml`文件。
2. 配置`env`变量为`debug`。
3. 初始化数据库，见`zhuque.sql`。
4. 使用`test`账号，密码`test`登录系统。

### 分支管理
develop 本地开发分支
release 测试环境使用
master 生产环境使用

### 部署流程
1. 在服务器指定位置下载源码。
2. 安装go环境。go可以交叉编译，但是由于sqlite的缘故，windows环境下并不能顺利的编译linux版本，所以最好还是在linux环境下编译。
3. 配置`conf.yaml`文件，参照`conf-sample.yaml`文件。
4. 项目目录中编译`go build`，第一次会安装依赖会慢一些。
5. `./zhuque`启动服务。

### 权限架构
该系统权限使用了基于角色的访问控制方法（RBAC）。是目前公认的解决大型企业的统一资源访问控制的有效方法。 其显著的两大特征是：1.减小授权管理的复杂性，降低管理开销。2.灵活地支持企业的安全策略，并对企业的变化有很大的伸缩性。

参考文档：

[官方原文](https://pm2.keymetrics.io/docs/usage/deployment/)

[PM2自动部署代码流程总结](https://segmentfault.com/a/1190000017310047)

[pm2 官方文档 学习笔记](https://my.oschina.net/u/4400196/blog/3283439)

[通过Github与PM2部署Node应用](https://zhuanlan.zhihu.com/p/20940096)

[deploy](https://github.com/Unitech/PM2/blob/0.14.7/ADVANCED_README.md#deployment-options)

### [简单部署](#simple-deploy)

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

### [Complete tutorial](#complete-tutorial)

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

### 应用选项

属性`apps`是包含了一组`Object`的数组。

选项名称|描述|类型|默认
---|---|---|---
`script`|启动脚本的路径，必填字段|`String`|
`name`|进程列表中的进程名称|`String`|没有扩展名的脚本文件名（`app.js` 对应的名字是`app`）
`cwd`|当前工作目录以启动进程|`String`|当前`shell`环境的`CWD`
`args`|传递给脚本的参数|`Array`,`String`|
`interpreter`|解释器绝对路径|`String`|`node`
`node_args`|传递给解释器的参数|`Array`,`String`|
`output`|`stdout`的文件路径（输出的行会追加到文件中）|`String`|`~/.pm2/logs/<app_name>-out.log`
`error`|`stderr`的文件路径（输出的行会追加到文件中）|`String`|`~/.pm2/logs/<app_name>-error.err`
`log`|组合`stdout`和`stderr`的文件路径（输出的行会追加到文件中)|`Boolean`,`String`|`/dev/null`
`disable_logs`|禁用日志存储|`Boolean`|
`log_type`|指定日志输出类型，可能的值为：`json`|`String`|
`log_date_format`|日志时间戳的格式，采用`moment.js`格式（例如`YYYY-MM-DD HH：mm Z`）|`String`|
`env`|指定要注入的环境变量|`Object`,`String`|
`^env_\S*$`|指定使用`--env <env_name>`时要注入的环境变量|`Object`,`String`|
`max_memory_restart`|如果进程内存超出这个指定的最大内存，会重新启动应用（格式：`[0-9]` ，`K`是`KB`, `M`是`MB`, `G`是`GB`, 默认是`B`）|`String`,`Number`|
`pid_file`|已启动进程的`pid`的文件路径，有`pm2`写入|`String`|`~/.pm2/pids/app_name-id.pid`
`restart_delay`|在重启崩溃应用的延迟时间，单位毫秒|`Number`|
`source_map_support`|启用或禁用源映射支持|`Boolean`|`true`
`disable_source_map_support`|启用或禁用源映射支持|`Boolean`|
`wait_ready`|如果为`true`，进程会等待`process.send（'ready'）`的`ready`事件|`Boolean`|
`instances`|指定集群模式下启动的实例数|`Number`|`1`
`kill_timeout`|指定`PM2`杀掉进程的超时时间，单位毫秒。`PM2`在发送`SIGKILL`信号之前，如果应用进程没用自己退出，`PM2`等待{kill_timeout}毫秒后会主动杀掉应用进程|`Number`|`1600`
`listen_timeout`|单位毫秒，`PM2`会监听应用是否`ready`，如果{listen_timeout}后没有`ready`会强制重载，则强制重载|`Number`|
`cron_restart`|一个`cron`模式来重启你的应用|`String`|
`merge_logs`|在集群模式下，将每种类型的日志合并到一个文件中（而不是每个集群都有单独的日志文件）|`Boolean`|
`vizion`|启用或禁用版本元数据（`vizion`库）|`Boolean`|`true`
`autorestart`|进程失败后启用或禁用自重启|`Boolean`|`true`
`watch`|启用或禁用观察模式|`Boolean`,`Array`,`String`|
`ignore_watch`|要忽略的路径列表（正则表达式）|`Array`,`String`|
`watch_options`|用作`chokidar`选项的对象（请参阅`chokidar`文档）|`Object`|
`min_uptime`|考虑应用启动的最小正常运行时间（格式为`[0-9]`, `h`小时，`m`分钟，`s`秒，默认为`ms`）|`Number`,`String`|`1000`
`max_restarts`|最多重启次数|`Number`|`16`
`exec_mode`|设置执行模式，可能的值为：`fork|cluster`|`String`|`fork`
`force`|即使脚本已经运行，也强制将其启动|`Boolean`|
`append_env_to_name`|将环境名称附加到应用名称|`Boolean`|
`post_update`|在从`Keymetrics`仪表板执行提取/升级操作之后执行的命令列表|`Array`|
`trace`|启用或禁用事务跟踪|`Boolean`|
`disable_trace`|启用或禁用事务跟踪|`Boolean`|`true`
`increment_var`|指定环境变量的名称以注入每个群集的增量|`String`|
`instance_var`|重命名`NODE_APP_INSTANCE`环境变量|`String`|`NODE_APP_INSTANCE`
`pmx`|开启或禁用`pmx`包装|`Boolean`|`true`
`automation`|开启或禁用`pmx`包装|`Boolean`|`true`
`treekill`|只`kill`主进程，不分离子进程|`Boolean`|`true`
`port`|注入PORT环境变量的快捷方式|`Number`|
`uid`|设置用户`ID`|`String`|当前用户的`uid`
`gid`|设置群组`ID`|`String`|当前用户的`uid`

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

