# POC RabbitMQ using Dead Letter Message (dlx) and TTL

## Running RabbitMQ Server through docker-compose

```sh
docker-compose up -d
```

## Create Exchange Maually

Create exchange through rabbit-admin-ui

- `e2`
- `e-dlx`

## Create Queue using Programmatically

```sh
go run cmd/consumer/main.go
go run cmd/consumer-dlx/main.go
```

Folloing above command will create quete and binding queue with exchange

## Resources

- https://dzone.com/articles/rabbitmq-consumer-retry-mechanism
- https://www.rabbitmq.com/dlx.html
- https://www.rabbitmq.com/ttl.html
