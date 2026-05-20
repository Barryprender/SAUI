FROM golang:1.25-alpine AS builder

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
