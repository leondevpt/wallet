FROM golang:alpine AS builder
RUN apk add --no-cache ca-certificates build-base libc6-compat make git bash
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make build



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
COPY --from=builder /app/bin ./bin/

VOLUME ["/app/config", "/app/logs", "/app/data"]
USER appuser
EXPOSE 50051

ENTRYPOINT ["./bin/app"]