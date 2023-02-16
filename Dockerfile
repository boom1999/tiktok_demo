FROM golang:1.19-alpine

WORKDIR /apps

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod tidy \
    && go build -o tiktok_demo main.go

EXPOSE 8080

CMD ["/apps/tiktok_demo"]