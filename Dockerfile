# Stage 1: сборка
FROM golang:1.23 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o avito_merch ./cmd/main.go

# Stage 2: минимальный образ
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/avito_merch .
COPY --from=builder /app/web ./web
EXPOSE 8080
CMD ["./avito_merch"]