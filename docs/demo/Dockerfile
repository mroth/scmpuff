
FROM golang:1.26-alpine AS builder
RUN go install github.com/mroth/scmpuff@latest

FROM ghcr.io/charmbracelet/vhs:latest
RUN apt-get update && apt install -y git
COPY --from=builder /go/bin/scmpuff /usr/local/bin/scmpuff
