# syntax=docker/dockerfile:1
FROM golang:1.22-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o ipa-renamer-watch main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/ipa-renamer-watch /app/ipa-renamer-watch
COPY --from=builder /app/ipa_renamer /app/ipa_renamer
RUN chmod +x /app/ipa-renamer-watch /app/ipa_renamer
ENV WATCH_DIR=/app/watched
ENV IPA_RENAMER=/app/ipa_renamer
ENV OUTPUT_DIR=/app/output
CMD ["/app/ipa-renamer-watch"]
