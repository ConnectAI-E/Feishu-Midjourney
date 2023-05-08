FROM golang:1.20 AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /lark

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build lark lark

EXPOSE 16008 

USER nonroot:nonroot

ENTRYPOINT ["/lark"]
