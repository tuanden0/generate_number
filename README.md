# Generate Number

A small service that generate random number using buf.build, protobuf, gRPC, gRPC gateway and validator v10.

You can get API prototype [here](proto/gen/go/v1/number/number.swagger.json) and paste it to [Swagger UI Editor](https://editor.swagger.io/)

## How to run
```bash
# Server
go run cmd\numberapis\v1\server\main.go

# Client
go run cmd\numberapis\v1\client\main.go
```

## How to test
```bash
curl --location --request POST '127.0.0.1:8000/v1/get_num' \
--header 'Content-Type: application/json' \
--data-raw '{
    "num_a": 10,
    "num_b": 50,
    "num_c": 100,
    "num_d": -10
}'
```
