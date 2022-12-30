# Go Redis Timestamp Based Generator

[Redis Notes](docs/REDISNOTES.md)

## Description

A Redis library for processing timestamp based elements based on an external timestamp update.

### Keywords

- **Latest timestamp**: A stored timestamp that is continuously incremented based on external updates. Can only increase.

- **Timestamp based element**: An element with an assiociated timestamp. This element should be processed only once the latest timestamp is updated to a value greater than the elements timestamp.

### Problems

- Must handle distributed processing.
- Elements must be handled at-least-once.

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

## Edge Case Scenarios

### Error during timestamp update flow

If there is an error during processing a timestamp update, the entire function should return an error. This function is safe to retry. Some timed elements might be re-processed.

### Timestamp update -> element creation -> timestamp updating processing

**Case 1**

1. Timestamp update.
2. Element is created with lower timestamp.
3. Timestamp update processing.

In this case, the element is processed as part of the `Add` method.

**Case 2**

1. Timestamp update.
2. Element is created with higher timestamp.
3. Timestamp update processing.

In this case, the new element is not processed as expected. The element is added to the list of elements safely, and can not be removed accidentially since processed elements are
removed by key (member name).

### Element creation -> timestamp update

**Case 1**

Potential issue:

1. New element is determined to be higher than latest timestamp.
2. Timestamp update is triggered with a timestamp higher than the previous element. New element is not added to the list yet. Timestamp update finishes processing without the previous element.
3. New element is added to list.

In this case, the new element is not handled, but the timestamp is updated to a higher value. However, this library handles this case by doing the following:

1. New element is determined to be higher than latest timestamp.
2. Timestamp update is triggered with a timestamp higher than the previous element. The command update timestamp is queued in Redis, but is behind the command to add the previous timed element to the set (since the commands are part of a Lua script).
3. New element is added to the list.
4. Timestamp update updates timestamp.
5. Timestamp update gets a list of all elements before this timestamp, which includes the new element from (1).
6. New element is processed.
