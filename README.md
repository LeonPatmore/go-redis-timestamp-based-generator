# Go Redis Timestamp Based Generator

[Redis Notes](docs/REDISNOTES.md)

## Example

[Basic Usage Example](cmd/time_element_example/main.go)

To run the example yourself:

`make run cmd=time_element_example`

<details>
    <summary>Expected Output</summary>

    Adding 3 new timed elements...
    Before timestamp update:
    [cfb76172-9829-4940-bf27-50924e75147c has timestamp 1 be4b1b4f-f0b4-49ce-9034-742b45aaf877 has timestamp 2 d76118d0-084d-47b4-a4bb-4fb7b76b1605 has timestamp 3]
    Pushing timestamp update for timestamp [ 2 ]. This should trigger handling of two elements.
    Handling element with ID cfb76172-9829-4940-bf27-50924e75147c
    Handling element with ID be4b1b4f-f0b4-49ce-9034-742b45aaf877
    After timestamp update:
    [d76118d0-084d-47b4-a4bb-4fb7b76b1605 has timestamp 3]
    Pushing element with timestamp [ 2 ]. This element should be handled instantly, and not added to the queue.
    Handling element with data [ 408b52dc-f312-4ac0-95f4-6ef480f56121 ] on ADD
    Handling element with ID 408b52dc-f312-4ac0-95f4-6ef480f56121
    After new element:
    [d76118d0-084d-47b4-a4bb-4fb7b76b1605 has timestamp 3]
    Pushing element with timestamp [ 3 ]. This element should not be handled instantly.
    After new element:
    [d76118d0-084d-47b4-a4bb-4fb7b76b1605 has timestamp 3 f121dc68-32bd-49d2-8e5b-338d2874f22d has timestamp 3]
    Pushing timestamp update for timestamp [ 3 ]. This should trigger handling of two elements.
    Handling element with ID d76118d0-084d-47b4-a4bb-4fb7b76b1605
    Handling element with ID f121dc68-32bd-49d2-8e5b-338d2874f22d
    There should now be zero elements left:
    []

</details>

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
