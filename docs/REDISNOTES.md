# Redis Notes

## Data Types

- String
- List (of strings)
- Set
- Sorted Set
- Hash

### Redis String

A Redis string is a simple data type for storing a sequence of bytes.
A Redis string can actually be all of the following:

- A string of characters.
- A byte array.
- A long number.
- An array.
- A block of memory.

**Implementation**

Uses a hashed location. Also stores the length of the string.

**Performance**

Set/get: O(1)
Setnx (Set if it doesn't already exist): O(1)
Substr/getrange/setrange: O(n)

### Redis Lists

A list of strings.

**Implementation**

Doubly linked list. Also stores the length of the list.

**Performance**

Operations at head or tail of list: O(1)
Operations on elements within the list: O(n)

### Redis Sets

Unordered collection of unique strings.

**Implementation**

Hash table.

**Performance**

Adding/removing/checking: O(1)
Getting all members: O(n)

### Redis Sorted Set

A set which members are ordered by an assiocated score. When the scores are equal for two members, they are ordered lexicographically.

**Implementation**

Skip list.

![alt text](skiplist.png)

In Redis's implementation, each node has a pointer back to the start of the list as-well, which allows one to traverse in reverse order of the score.

**Performance**

Most sorted set operations are O(log(n)).

### Redis Hash

A Redis hash is a collection of key/value pairs, used to represent basic objects.

**Implementation**

Hash table.

**Performance**

Set/get: O(1)
Get all keys/values: O(n)
