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
│    ├─ ffmpeg
│    ├─ jwt
│       ├─ Auth.go
│    ├─ redis
│    └─ others(etc., ssh connect)
├─ repostitory（dao层）
│    ├─ comment.go
│    ├─ follow.go
│    ├─ init.go
│    ├─ like.go
│    ├─ user.go
│    ├─ video.go
│    └─ others(etc., publish && message)
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
│    └─ others(etc., publish && message)
├─ util（工具类）
├─ docker-compose.yml（多容器管理）
├─ Dockerfile（容器指令）
├─ go.mod
├─ go.sum
├─ LICENSE
├─ main.go（主程序入口）
├─ README.md
└─  wait-for-it.sh
```

[license-shield]: https://img.shields.io/github/license/mrxuexi/tiktok.svg?style=flat-square

[license-url]: https://github.com/boom1999/tiktok_demo/blob/master/LICENSE

---
- Step 1. Fork this repository
- Step 2. Permission
    ```shell
    sudo chmod 777 wait-for-it.sh
    ```
- Step 3. Change your configs in `./config/config.yaml`
- Step 4. Install `docker` and `docker-compose`
- Step 5. `docker-compose up`
- ---