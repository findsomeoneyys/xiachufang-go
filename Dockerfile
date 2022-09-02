# FROM 基于 alpine:latest
FROM alpine:latest AS base
# EXPOSE 设置端口映射
EXPOSE 8080/tcp
# RUN 设置代理镜像
RUN echo -e http://mirrors.ustc.edu.cn/alpine/v3.13/main/ > /etc/apk/repositories
# RUN 设置 Asia/Shanghai 时区
RUN apk --no-cache add tzdata  && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone


FROM golang:1.18-alpine AS build
# ENV 设置环境变量
ENV GOPATH=/opt/repo
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . $GOPATH/src/github.com/findsomeoneyys/xiachufang-api
RUN cd $GOPATH/src/github.com/findsomeoneyys/xiachufang-api && go build .


FROM base AS final
# COPY 源路径 目标路径 从镜像中 COPY
COPY --from=builder /opt/repo/src/github.com/findsomeoneyys/xiachufang-api /opt
# WORKDIR 设置工作目录
WORKDIR /opt
# CMD 设置启动命令
CMD ["./xiachufang-api"]