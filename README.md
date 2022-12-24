# Go Redis Timestamp Based Generator

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

## Resources

- https://redis.io/docs/manual/patterns/distributed-locks/
- https://kafka.apache.org/documentation/#producerconfigs
