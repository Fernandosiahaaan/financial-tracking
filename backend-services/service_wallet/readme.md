# Service Wallet

## How to Run Local

```
cd ./financial-tracking/backend-services/service_wallet
go run /app/main.go
```


## How to Build

```
cd ./financial-tracking/backend-services/service_wallet
docker build -t service-wallet:v0.0.0 .
docker run --env-file .env --name service-wallet-cont service-wallet:v0.0.0
```