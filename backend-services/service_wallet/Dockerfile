# ============================
# Tahap Builder
# ============================
FROM golang:1.20.3 AS builder

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p build && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/service-wallet-app ./app/

# ============================
# Tahap Produksi
# ============================
FROM debian:bullseye-slim

WORKDIR /app

COPY .env .env
COPY --from=builder /app/build ./build

CMD ["./build/service-wallet-app"]
