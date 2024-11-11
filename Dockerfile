FROM golang:1.23.3 AS builder

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY ./main.go  ./

ARG CGO_ENABLED=0
ARG GOOS=linux
RUN go build \
    -o /go/bin/aws-config-compliance-prometheus-exporter \
    -ldflags '-s -w'

FROM alpine:3.20.3 AS runner

COPY --from=builder /go/bin/aws-config-compliance-prometheus-exporter /app/aws-config-compliance-prometheus-exporter

RUN adduser -D -S -H exporter

USER exporter

ENTRYPOINT ["/app/aws-config-compliance-prometheus-exporter"]
