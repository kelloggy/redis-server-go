# redis-server-go

This repository contains a simple Redis server implemented in Go, providing basic functionality for GET, SET, HSET, and HGET operations.

## Features
GET: Retrieve the value associated with a key.
SET: Set the value for a given key.
HSET: Set the field in the hash stored at the specified key to the provided value.
HGET: Retrieve the value of a field in the hash stored at the specified key.

## Usage
```
# Set a key-value pair
$ SET key1 value1

# Get the value associated with a key
$ GET key1

# Set a field in a hash
$ HSET hash1 field1 value1

# Get the value of a field in a hash
$ HGET hash1 field1
```
