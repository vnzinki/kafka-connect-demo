FROM golang:1.17-alpine AS builder

WORKDIR /build

ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV GOPROXY=https://proxy.golang.org

RUN apk --no-cache add \
    libc-dev gcc bash git \
    && rm -rf /var/cache/apk/*

COPY . /build
WORKDIR /build

RUN go mod download && \
    go build -o ./dist/healthz ./healthz.go

FROM confluentinc/cp-kafka-connect:6.2.4
# FROM bitnami/kafka:2.5.0
RUN confluent-hub install --no-prompt debezium/debezium-connector-mysql:latest &&\
    confluent-hub install --no-prompt debezium/debezium-connector-postgresql:latest &&\
    confluent-hub install --no-prompt debezium/debezium-connector-sqlserver:latest &&\
    confluent-hub install --no-prompt confluentinc/kafka-connect-jdbc:latest &&\
    confluent-hub install --no-prompt confluentinc/kafka-connect-elasticsearch:latest

COPY --from=builder /build/dist/healthz /home/appuser/healthz
