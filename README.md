## 简介

零基础使用 Gin 创建一个 web 项目.附带教程

本项目经过我公司线上验证，完善了日志监控系统，提供了通过golang并发特性提高性能的demo

## 项目重点

- service/list.go中详细展示了 golang 协程 通道  waitgroup 锁 组合用法，来实现数据并发处理的逻辑
- service/list.go中展示了 临时对象池和 本地缓存 优化性能的方法
- logrus日志 结合 sentry 进行时事报警
- logrus日志 实时写入es



## 版本

- 0.1.0 项目初始化
- v0.2.0 读取配置
- v0.3.0 记录日志
- v0.4.0 连接数据库
- v0.5.0 定义错误码
- v0.6.0 读取请求返回响应
- v0.7.0 添加核心逻辑, 用户数据的 CRUD
- v0.8.0 增加中间件
- v0.9.0 添加 JWT 认证
- v0.10.0 添加 HTTPS
- v0.11.0 添加 Makefile
- v0.12.0 添加版本信息
- v0.13.0 添加启动脚本
- v0.14.0 添加 Nginx 配置
- v0.15.0 添加测试的例子
- v0.16.0 添加 swagger 文档
- v0.17.0 完善部署方式
- v0.18.0 增加bigcache本地缓存，临时对象池sync.pool(service/list.go中有示例)
## 运行

假设: 在项目根目录下运行命令

方式一: 在 docker 中运行 mysql, 本地启动服务器

```bash
# 后台启动 mysql 服务器
docker-compose up -d mysql
# 运行服务器
go run ./
```

方式二: 在 docker 中运行 mysql, 本地编译二进制文件, 直接启动

```bash
# 后台启动 mysql 服务器
docker-compose up -d mysql
# 编译, 应该会在当前目录下生成一个叫做 web 的二进制文件
make build
# 运行
web
```

方式三: 在 docker 中运行 mysql, 使用 systemd 接管服务

额外要求: 目录应该是 /home/go_web/, 否则需要更改配置中的路径

使用的配置路径是 /home/go_web/conf/config_abs.yaml

```bash
# 后台启动 mysql 服务器
docker-compose up -d mysql
# 编译, 应该会在当前目录下生成一个叫做 web 的二进制文件
make build
# 复制文件到 systemd 的配置文件夹
cp conf/go_web.service /etc/systemd/system/
# 启动
systemctl start go_web
# 查看状态
systemctl status go_web
# 停止
systemctl stop go_web
```

方式四: 完全运行在 docker-compose 中

```bash
# 后台启动 mysql 服务器
docker-compose up -d mysql
# 预先构建 app
docker-compose build app
# 启动, 注意这里运行了 3 个 app
docker-compose up --scale app=3 nginx
```
