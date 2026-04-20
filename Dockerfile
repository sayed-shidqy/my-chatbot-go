# Ubah dari 1.21 ke 1.22 atau 1.23
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod init chat-bot-go || true
RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
EXPOSE 8080
CMD ["./main"]
