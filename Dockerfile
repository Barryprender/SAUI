# Pinned to match the `go` line in go.mod and the CI toolchain. The floor is
# 1.25.14: every earlier 1.25 patch ships standard-library vulnerabilities that
# govulncheck reports as reachable from this code. Raise all three together.
FROM golang:1.25.14-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o saui .

FROM alpine:3.21

RUN addgroup -S saui && adduser -S saui -G saui

WORKDIR /app

COPY --from=builder /app/saui ./
COPY --from=builder /app/static ./static

RUN mkdir -p /data && chown saui:saui /data

USER saui

EXPOSE 8080

CMD ["./saui"]
