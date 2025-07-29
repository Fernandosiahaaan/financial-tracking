# Service User

## How to Run Local

```
cd ./financial-tracking/backend-services/service_user
go run /app/main.go
```


## How to Build

```
cd ./financial-tracking/backend-services/service_user
docker build -t service_user:v0.0.0 .
docker run --env-file .env --name service-user-cont service_user:v0.0.0
```