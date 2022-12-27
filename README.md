# Go Redis Timestamp Based Generator

[Redis Notes](docs/REDISNOTES.md)

## Requirements

- Docker
- Go

## Redis Version

7

## Atmoic

- Redis is single threaded, so commands can not happen at the same time.

## Consistency

- A Redis cluster is not strongly consistent. This is because write acks from replicas are all
  async.

## Locking

Using https://github.com/go-redsync/redsync.

## Resources

- https://redis.io/docs/manual/patterns/distributed-locks/
- https://kafka.apache.org/documentation/#producerconfigs
