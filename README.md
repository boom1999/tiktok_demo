<!-- PROJECT SHIELDS -->


## Tiktok Demo

### *Updating*

![GitHub Repo stars](https://img.shields.io/github/stars/boom1999/tiktok_demo??style=plastic)
![GitHub watchers](https://img.shields.io/github/watchers/boom1999/tiktok_demo??style=plastic)
![GitHub forks](https://img.shields.io/github/forks/boom1999/tiktok_demo??style=plastic)
![GitHub contributors](https://img.shields.io/github/contributors/boom1999/tiktok_demo??style=plastic)
[![MIT License][license-shield]][license-url]

![golang](https://img.shields.io/badge/golang-1.19-blue)
![golang](https://img.shields.io/badge/gorm-1.8.2-red)
![golang](https://img.shields.io/badge/gorm-1.24.5-green)
![golang](https://img.shields.io/badge/viper-1.15.0-orange")

---
```
├─ config（配置文件信息以及Viper管理文件）
│    ├─ config.go
│    └─ config.yml
├─ controller（处理客户端请求的控制层）
│    ├─ comment.go
│    ├─ follow.go
│    ├─ like.go
│    ├─ video.go
│    ├─ user.go
│    └─ others(etc., publish && message)
├─ middleware（中间件）
│    ├─ jwt
│    │    └─ Auth.go
│    ├─ minio
│    │    ├─ init.go
│    │    └─ utils.go
│    ├─ rabbitmq
│    │    ├─ commentMQ.go
│    │    ├─ followMQ.go
│    │    ├─ likeMQ.go
│    │    └─ rabbitMQ.go
│    └─ redis
│         └─ cache.go
├─ repostitory（dao层）
│    ├─ comment.go
│    ├─ follow.go
│    ├─ init.go
│    ├─ like.go
│    ├─ user.go
│    ├─ video.go
│    └─ message.go
├─ routes（路由层）
│    ├─ comment.go
│    ├─ favorite.go
│    ├─ message.go
│    ├─ publish.go
│    ├─ relation.go
│    ├─ routes.go
│    └─ user.go
├─ service（业务逻辑层）
│    ├─ comment.go
│    ├─ follow.go
│    ├─ like.go
│    ├─ user.go
│    ├─ video.go
│    └─ messgae.go
├─ util（工具类）
│    └─ util.go
├─ .gitignore
├─ docker-compose.yml（多容器管理）
├─ Dockerfile（容器指令）
├─ go.mod
├─ go.sum
├─ LICENSE
├─ main.go（主程序入口）
├─ README.md
└─ wait-for-it.sh
```

[license-shield]: https://img.shields.io/github/license/mrxuexi/tiktok.svg?style=flat-square

[license-url]: https://github.com/boom1999/tiktok_demo/blob/master/LICENSE

---
> In order to make data portable and reusable, we use **volumes** to mount docker data.
> 
> Before that, please open each **port** of the corresponding service.

- Step 1. Fork this repository
  ``` shell
  git clone -b master git@github.com:boom1999/tiktok_demo.git
  ```
- Step 2. Change your configs in `./config/config.yaml`, especially for `host`
- Step 3. Install `Docker` and `docker-compose`
  - Install `Docker`: 
    ``` shell
    curl -fsSL https://get.docker.com | bash -s docker`, also you can add `--mirror Aliyun
    ```
  - Start docker service: `systemctl start docker`
  - Download `docker-compose`
    ```shell
    curl -L "https://github.com/docker/compose/releases/download/v2.2.3/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    ```
  - Read and write permission: `chmod -R 777 /usr/local/bin/docker-compose`
- Step 4. File permission 
  - Wait for dependent services to start
    ```shell
      chmod 777 wait-for-it.sh
      ```
    tips: maybe need to change its format to unix(if not), `:set ff=unix`
  - For rabbitMQ log file
    - `mkdir -p tiktokCache/rabbitMQ/log`
    - Give permissions to log and all subdirectories
      ```shell
      cd tiktokCache/rabbitMQ
      chmod -R 777 log/
      ```
  - For redis conf file
    - `mkdir -p tiktokCache/redis/conf`
    - download `redis.conf` from [Redis.conf](https://redis.io/docs/management/config/) and copy it to conf directory
    - `mkdir -p tiktokCache/redis/logs`
    - Give permissions to log and all subdirectories
      ```shell
      cd tiktokCache/redis
      chmod -R 777 logs/
      ```
- Step 5.RUN: 
  ```shell
  docker-compose up
  ```
**Tips:**

- Delete all the containers
```shell
 docker rm $(docker ps -a -q)
```
- Remove `none` image
```shell
docker image prune
```

> If you are deploying on a remote **ECS** instead of a **virtual machine** and want to connect through tools such as _Navicat_,
> please make sure that mysql has enabled the remote connection permission for the _user_ or _root_ and FirewallD port 3306 
> (if not enabled, it will not affect data reading and writing).