#!/usr/bin/env bash

docker run -it -v "$(pwd):/app" -w /app golang:1.25.5 \
  sh -c 'go mod tidy && \
  GOOS=linux GOARCH=amd64 go build -o vmStarter-linux-amd64 && \
  GOOS=darwin GOARCH=arm64 go build -o vmStarter-darwin-arm64'

REMOTE="user@127.0.0.1"
APP="/opt/tools/vmStarter"
scp vmStarter-linux-amd64 ${REMOTE}:${APP}
ssh ${REMOTE} chmod +x ${APP}
