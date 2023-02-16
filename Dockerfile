FROM golang:1.19-alpine as builder

WORKDIR /apps

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod tidy \
    && go build -o tiktok_demo main.go

FROM alpine:3.14

WORKDIR /root/
RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache ffmpeg
COPY --from=builder /apps/tiktok_demo ./tiktok_demo_ffmpeg

EXPOSE 8080

CMD ["./wait-for-it.sh", "mysql:3306","rabbitmq:5672","-t 45","--", "./tiktok_demo_ffmpeg" ]