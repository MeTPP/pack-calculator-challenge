FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o pack-calculator ./cmd/main.go

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/pack-calculator .
COPY --from=builder /app/templates/ ./templates/

RUN addgroup -S norootgroup && adduser -S noroot -G norootgroup
RUN chown -R noroot:norootgroup /app
USER noroot

EXPOSE 8080

CMD ["./pack-calculator"]
