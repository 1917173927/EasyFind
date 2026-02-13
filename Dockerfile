FROM golang:1.25-alpine AS builder

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

FROM alpine:latest

RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/config ./config
# Create uploads directory
RUN mkdir -p uploads

EXPOSE 8080

CMD ["./main"]
