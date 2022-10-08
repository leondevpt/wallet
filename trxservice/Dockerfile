FROM golang:alpine3.15 as builder
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
# 添加cgo 相关的库
 && apk update && apk upgrade && apk add ca-certificates  build-base libc6-compat
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o app


# Run container
FROM alpine3.15
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk --no-cache add  tzdata bash vim ca-certificates && \
    update-ca-certificates && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# Create a group and user
RUN addgroup appuser --gid 1000 && \
    adduser appuser --ingroup appuser --uid 1000 --home /home/appuser --disabled-password && \
    mkdir -p /app/bin && mkdir /app/config && mkdir /app/logs && mkdir /app/data  && chown appuser:appuser /app

WORKDIR /app
COPY --from=builder /app/app ./bin/

VOLUME ["/app/config", "/app/logs", "/app/data"]
USER appuser
EXPOSE 8080

ENTRYPOINT ["./bin/app"]
