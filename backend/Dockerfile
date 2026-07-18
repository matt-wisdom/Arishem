FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN CGO_ENABLED=1 GOOS=linux go build -o /server ./cmd/server

FROM alpine:3.19

RUN apk add --no-cache ca-certificates docker-cli python3 py3-pip

WORKDIR /app

COPY --from=builder /server .

COPY pyproject.toml arishem ./
RUN pip install --no-deps --break-system-packages -e /app

EXPOSE 8080

CMD ["./server"]