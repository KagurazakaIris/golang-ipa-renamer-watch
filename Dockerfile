# syntax=docker/dockerfile:1.7
FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o ipa-renamer .

FROM alpine:3.20
WORKDIR /app
LABEL org.opencontainers.image.description="跨平台 Go 目录监听自动重命名 IPA 工具，命名规则：原文件名@CFBundleIdentifier.ipa，支持 Docker/Compose/GHCR 多平台镜像。"
COPY --from=builder /app/ipa-renamer /app/ipa-renamer
RUN chmod +x /app/ipa-renamer
ENV WATCH_DIR=/app/watched
ENV OUTPUT_DIR=/app/output
ENV TEMP_DIR=/app/temp
CMD ["/app/ipa-renamer"]
