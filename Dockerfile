# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o fo-sentinel .

# Runtime stage
FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/fo-sentinel .
COPY --from=builder /app/manifest ./manifest

ENV TZ=Asia/Shanghai

EXPOSE 8000

CMD ["./fo-sentinel"]
