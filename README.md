# Go Redis Timestamp Based Generator

[Redis Notes](docs/REDISNOTES.md)

## Requirements

- Docker
- Go

### Redis Version

7

## Testing

Ensure you have set the following two env vars:

- `REDIS_HOST`: Host of Redis server.
- `REDIS_PORT`: Port of Redis server.

Then run `make test`.

## Resources

- https://redis.io/docs/manual/patterns/distributed-locks/
- https://kafka.apache.org/documentation/#producerconfigs
