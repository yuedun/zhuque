# 朱雀-nodejs发布系统
朱雀发布系统是专门为nodejs发布而开发的系统，使用go语言开发，前端使用了layui mini框架，简单易上手。
朱雀发布系统前后端分离，但又是一体的，无需部署两套服务。
理论上朱雀发布系统可以发布其他语言应用程序，但其与nodejs更加相得益彰。
原因是朱雀发布系统依赖于PM2，PM2大家都知道，与nodejs几乎是绝配，所以有nodejs的地方很大可能有PM2，所以使用朱雀发布系统就显得更加简单了，无需专门安装PM2。

数据库使用sqlite，无需单独安装和配置。

# 使用框架文档
[前端框架](http://layuimini.99php.cn/docs/index.html)
[前端框架](http://layuimini.99php.cn/onepage/v2/index.html)
[前端框架](https://www.layui.com/doc/)
[Font Awesome图标库](https://fontawesome.dashgame.com/)
[后端框架gorm](http://gorm.book.jasperxu.com/)
[后端框架gin](https://github.com/gin-gonic/gin#using-middleware)

# 本地开发
## 首次运行
1. 复制`conf-sample.yaml`文件为`conf.yaml`文件。
2. 配置`env`变量为`debug`。
3. 修改`\zhuque\pkg\user\handler.go`,`Init`函数的返回数据`menuInfo`为注释代码，因为第一次运行系统没有系统数据，需要模拟数据。
4. 使用`test`账号，密码test登录系统，增加角色数据，用户数据，分配权限。然后可以使用新用户登录。

## 部署流程
1. 在服务器指定位置下载源码。
2. 安装go环境。go可以交叉编译，但是由于sqlite的缘故，windows环境下并不能顺利的编译linux版本，所以最好还是在linux环境下编译。
3. 配置`conf.yaml`文件，参照`conf-sample.yaml`文件。
4. 项目目录中编译`go build`，第一次会安装依赖会慢一些。
5. `./zhuque`启动服务。

## 权限架构
该系统权限使用了基于角色的访问控制方法（RBAC）。是目前公认的解决大型企业的统一资源访问控制的有效方法。 其显著的两大特征是：1.减小授权管理的复杂性，降低管理开销。2.灵活地支持企业的安全策略，并对企业的变化有很大的伸缩性。

参考文档：

[官方原文](https://pm2.keymetrics.io/docs/usage/deployment/)

[PM2自动部署代码流程总结](https://segmentfault.com/a/1190000017310047)

[pm2 官方文档 学习笔记](https://my.oschina.net/u/4400196/blog/3283439)

[通过Github与PM2部署Node应用](https://zhuanlan.zhihu.com/p/20940096)

[deploy](https://github.com/Unitech/PM2/blob/0.14.7/ADVANCED_README.md#deployment-options)

## [Getting started](#getting-started)

PM2 embeds a simple and powerful deployment system with revision tracing.

Please read the [Considerations to use PM2 deploy](#considerations).

## [简单部署](#simple-deploy)

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

## [Complete tutorial](#complete-tutorial)

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

## [Deployment options](#deployment-options)

Display deploy help via `pm2 deploy help`:

```
pm2 deploy <configuration_file> <environment> <command>

  Commands:
    setup                run remote setup commands
    update               update deploy to the latest release
    revert [n]           revert to [n]th last deployment or 1
    curr[ent]            output current release commit
    prev[ious]           output previous release commit
    exec|run <cmd>       execute the given <cmd>
    list                 list previous deploy commits
    [ref]                deploy to [ref], the "ref" setting, or latest tag
```

## [Use different set of env variables](#use-different-set-of-env-variables)

In the `post-deploy` attribute, you may have noticed the command `pm2 startOrRestart ecosystem.json --env production`. The `--env <environment_name>` allows to inject different sets of environment variables.

Read more [here](http://pm2.keymetrics.io/docs/usage/application-declaration/#switching-to-different-environments).

## [Related Commands](#related-commands)

```
pm2 startOrRestart all.json            # Invoke restart on all apps in JSON
pm2 startOrReload all.json             # Invoke reload
```

## [Multi host deployment](#multi-host-deployment)

To deploy to multiple hosts in the same time, you just have to declare each host in an array under the attribute `host`.
```
{
  [...]
  "deploy" : {
    "production" : {
      "user" : "node",
      // Multi host in a js array
      "host" : ["212.83.163.1", "212.83.163.2", "212.83.163.3"],
      "ref"  : "origin/master",
      "repo" : "git@github.com:repo.git",
      "path" : "/var/www/production",
      "pre-setup" : "echo 'commands or local script path to be run on the host before the setup process starts'",
      "post-setup": "echo 'commands or a script path to be run on the host after cloning the repo'",
      "post-deploy" : "pm2 startOrRestart ecosystem.json --env production",
      "pre-deploy-local" : "echo 'This is a local executed command'"
    }
  [...]
}
```

## [Using SSH keys](#using-ssh-keys)

You just have to add the “key” attribute with path to the public key, see below example :
```
    "production" : {
      "key"  : "/path/to/some.pem", // path to the public key to authenticate
      "user" : "node",              // user used to authenticate
      "host" : "212.83.163.1",      // where to connect
      "ref"  : "origin/master",
      "repo" : "git@github.com:repo.git",
      "path" : "/var/www/production",
      "post-deploy" : "pm2 startOrRestart ecosystem.json --env production"
    },
```
## [Force deployment](#force-deployment)

You may get this message:

```
--> Deploying to dev environment
--> on host 192.168.1.XX

  push your changes before deploying

Deploy failed
```
That means that you have changes in your local system that aren’t pushed inside your git repository, and since the deploy script get the update via `git pull` they will not be on your server. If you want to deploy without pushing any data, you can append the `--force` option:

`pm2 deploy ecosystem.json production --force`

## [Considerations](#considerations)

*   You can use the option `--force` to skip local change detection
*   You might want to commit your node_modules folder ([#622](https://github.com/Unitech/pm2/issues/622)) or add the `npm install` command to the `post-deploy` section: `"post-deploy" : "npm install && pm2 startOrRestart ecosystem.json --env production"`
*   Verify that your remote server has the permission to git clone the repository
*   You can declare specific environment variables depending on the environment you want to deploy the code to. For instance to declare variables for the production environment, add “env_production”: {} and declare the variables.
*   By default, PM2 will use `ecosystem.json`. So you can skip the <configuration_file> options if this is the case</configuration_file>
*   You can embed the “apps” &amp; “deploy” section in the package.json
*   It deploys your code via ssh, you don’t need any dependencies
*   Processes are initialized / started automatically depending on the application name in `ecosystem.json`
*   PM2-deploy repository can be found here: [pm2-deploy](https://github.com/Unitech/pm2-deploy)
*   **WINDOWS** : see point below (at the end)

## [Troubleshooting](#troubleshooting)

##### SSH clone errors

In most cases, these errors will be caused by `pm2` not having the correct keys to clone your repository. You need to verify at every step that the keys are available.

**Step 1**
If you are certain your keys are correctly working, first try running `git clone your_repo.git` on the target server. If it succeeds, move onto the next steps. If it failed, make sure your keys are stored both on the server and on your git account.

**Step 2**
By default `ssh-copy-id` copies the default identiy, usually named `id_rsa`. If that is not the appropriate key:

`ssh-copy-id -i path/to/my/key your_username@server.com`

This adds your public key to the `~/.ssh/authorized_keys` file.

**Step 3**
If you get the following error:
```
--> Deploying to production environment
--> on host mysite.com
  ○ hook pre-setup
  ○ running setup
  ○ cloning git@github.com:user/repo.git
Cloning into '/var/www/app/source'...
Permission denied (publickey).
fatal: Could not read from remote repository.

Please make sure you have the correct access rights and that the repository exists.

**Failed to clone**

Deploy failed
```

…you may want to create a ssh config file. This is a sure way to ensure that the correct ssh keys are used for any given repository you’re trying to clone. See [this example](https://gist.github.com/Protosac/c3fb459b1a942f161f23556f61a67d66):

```
# ~/.ssh/config
Host alias
    HostName myserver.com
    User username
    IdentityFile ~/.ssh/mykey
# Usage: `ssh alias`
# Alternative: `ssh -i ~/.ssh/mykey username@myserver.com`

Host deployment
    HostName github.com
    User username
    IdentityFile ~/.ssh/github_rsa
# Usage:
# git@deployment:username/anyrepo.git
# This is for cloning any repo that uses that IdentityFile. This is a good way to make sure that your remote cloning commands use the appropriate key
```

## [Windows Consideration](#windows-consideration)

To run the deploy script under Windows, you need to use a unix shell like bash, so we recommend to install either [Git bash](https://git-scm.com/download/win), [Babun](http://babun.github.io/) or  [Cygwin](https://cygwin.com/install.html)

## [Contributing](#contributing)

The module is [https://github.com/Unitech/pm2-deploy](https://github.com/Unitech/pm2-deploy)
Feel free to PR for any changes or fix.

## 生态系统文件参考

生态系统文件的目的是收集应用所有的配置选项和环境变量。

它是一个`javascript`文件，`exports`一个包含所有配置选项的`object`。这个`object`有两个属性：
- `apps`, `Array` 一组应用的配置
- `deploy`, `Object ` 部署配置选项

```javascript
module.exports = {
  apps: [{}, {}],
  deploy: {}
}
```
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

属性`deploy`是一个`Object`，每个对象属性定义了一个环境的部署配置选项。

结构：
```javascript
module.exports = {
  apps: [{}, {}],
  deploy: {
    production: {},
    staging: {},
    development: {}
  }
}
```

环境的部署配置选项：

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

# SCP发布
scp发布是较为常用的发布方式，但朱雀使用的是rsync发布，主要是利用其增量同步功能，加速代码同步。

## 发布流程
推荐的做法是在发布机上拉代码，编译。
同步代码到远程应用服务器，重启。
也就是说远程应用服务只需重启即可，不需要做编译操作，在发布机上编译即可。

scp发布模式的配置文件和pm2一样，主要用到的三个参数是：
- build 编译命令，可选。如：npm run build
- rsyncArgs rsync参数，可选。如：--exclude node_modules
- post-deploy 在应用服务器上执行的命令，如：重启服务 pm2 restart app.js