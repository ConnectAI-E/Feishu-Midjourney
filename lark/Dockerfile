FROM golang:1.18 as golang

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOPROXY=https://goproxy.cn,https://goproxy.io,direct

WORKDIR /build

COPY . /build
ADD config.example.yaml /build/config.yaml

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -o lark_midjourney

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache bash
COPY --from=golang /build/lark_midjourney /app
COPY --from=golang /build/config.yaml /app
EXPOSE 16008
ENTRYPOINT ["/app/lark_midjourney"]
