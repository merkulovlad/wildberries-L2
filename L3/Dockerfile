FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/notify-service ./cmd/notify-service

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

USER app
WORKDIR /app

COPY --from=builder /bin/notify-service /app/notify-service

EXPOSE 8080

ENTRYPOINT ["/app/notify-service"]
