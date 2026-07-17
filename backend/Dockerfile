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

# Bake arishem source into container
COPY --from=builder /app/arishem /app/arishem
COPY --from=builder /app/pyproject.toml /app/
COPY --from=builder /app/arishem.config.toml /app/

EXPOSE 8080

# Copy source to shared tmp on startup
COPY --from=builder /app/arishem /tmp/arishem/arishem-source
COPY --from=builder /app/pyproject.toml /tmp/arishem/arishem-source/
COPY --from=builder /app/arishem.config.toml /tmp/arishem/arishem-source/

CMD ["./server"]