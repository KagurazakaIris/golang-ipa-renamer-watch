# syntax=docker/dockerfile:1.7
FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o ipa-renamer main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/ipa-renamer /app/ipa-renamer
RUN chmod +x /app/ipa-renamer
ENV WATCH_DIR=/app/watched
ENV OUTPUT_DIR=/app/output
ENV TEMPLATE=$raw@$CFBundleIdentifier
ENV TEMP_DIR=/app/temp
CMD ["/app/ipa-renamer"]
