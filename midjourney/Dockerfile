FROM golang:1.18 as golang

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build
COPY . /build
ADD config.example.yaml /build/config.yaml

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -o discord_midjourney

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache bash
COPY --from=golang /build/discord_midjourney /app
COPY --from=golang /build/config.yaml /app
EXPOSE 16007
ENTRYPOINT ["/app/discord_midjourney"]
