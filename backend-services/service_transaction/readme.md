# Service Category
this repo is a service category of finacial tracking

# How to Run Local
```
mvn clean install
mvn spring-boot:run
```

## How to Build

```
cd ./financial-tracking/backend-services/service_transaction
docker build -t service-transaction:v0.0.0 .
docker run --env-file .env --name service-transaction-cont service-transaction:v0.0.0
```