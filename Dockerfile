FROM golang:1.18

MAINTAINER zibesun <979271618@qq.com>

WORKDIR /build

ADD . ./

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn

RUN go mod download

RUN go build ./main.go

EXPOSE 8000

ENTRYPOINT ["./main"]