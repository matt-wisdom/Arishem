FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /server ./cmd/server

FROM alpine:3.19

RUN apk add --no-cache ca-certificates docker-cli python3 py3-pip

WORKDIR /app

COPY --from=builder /server .

# Arishem source will be provided at runtime via volume mount
# See docker-compose.yml volumes section

EXPOSE 8080

CMD ["./server"]