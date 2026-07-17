FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /server ./cmd/server

# Create source tarball for runtime
RUN tar -czf /arishem-source.tar.gz arishem pyproject.toml arishem.config.toml

FROM alpine:3.19

RUN apk add --no-cache ca-certificates docker-cli python3 py3-pip

WORKDIR /app

COPY --from=builder /server .

# Extract arishem source
COPY --from=builder /arishem-source.tar.gz /tmp/
RUN mkdir -p /arishem && tar -xzf /tmp/arishem-source.tar.gz -C /arishem && rm /arishem-source.tar.gz

EXPOSE 8080

CMD ["./server"]