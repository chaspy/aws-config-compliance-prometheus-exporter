FROM golang:1.20.5 AS builder

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY ./main.go  ./

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build \
    -o /go/bin/aws-config-compliance-prometheus-exporter \
    -ldflags '-s -w'

FROM alpine:3.18.2 AS runner

COPY --from=builder /go/bin/aws-config-compliance-prometheus-exporter /app/aws-config-compliance-prometheus-exporter

RUN adduser -D -S -H exporter

USER exporter

ENTRYPOINT ["/app/aws-config-compliance-prometheus-exporter"]
