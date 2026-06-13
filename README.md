# go-redis

A Redis-compatible in-memory key-value server written in Go from scratch. Implements the Redis Serialization Protocol (RESP) with an epoll-based async TCP server, TTL/expiry, multiple eviction policies, and AOF persistence — all with zero external dependencies beyond `golang.org/x/sys`.

---

## Features

- **RESP Protocol** — Full encode/decode for Simple Strings, Errors, Integers, Bulk Strings, and Arrays
- **Async TCP Server** — Linux epoll-based non-blocking I/O supporting up to 20,000 concurrent clients
- **Sync TCP Server** — Simple Go `net` blocking server (alternative implementation)
- **TTL & Expiry** — Per-key expiration with Redis-style active expiry sampling
- **Eviction Policies** — `simple-first`, `allkeys-random`, `allkeys-lru`
- **AOF Persistence** — Append-Only File dump on `BGREWRITEAOF` or graceful shutdown
- **Type/Encoding System** — INT, EMBSTR, RAW string encodings with type-safe assertions
- **Pipelining** — Decode multiple commands from a single read
- **Graceful Shutdown** — Catches SIGTERM/SIGINT, persists data, and exits cleanly
- **Keyspace Stats** — Track key counts via the `INFO` command

### Supported Commands

| Command        | Description                                        |
| -------------- | -------------------------------------------------- |
| `PING`         | Returns `PONG` (or echoes a provided argument)     |
| `SET`          | Store a key-value pair (supports `EX` in seconds)  |
| `GET`          | Retrieve a value by key                            |
| `TTL`          | Get remaining time-to-live in seconds              |
| `DEL`          | Delete one or more keys                            |
| `EXPIRE`       | Set a TTL on an existing key                       |
| `INCR`         | Atomically increment an integer-valued string      |
| `INFO`         | Return keyspace statistics                         |
| `BGREWRITEAOF` | Write all keys to the AOF file                     |
| `LRU`          | Trigger manual allkeys-LRU eviction                |
| `CLIENT`       | Stub (returns `+OK`)                               |
| `LATENCY`      | Stub (returns empty array)                         |
| `SLEEP`        | Block for N seconds (debug utility)                |

---

## Project Structure

```
go-redis/
├── main.go                      # Entry point — starts server & signal handler
├── go.mod / go.sum              # Go module definition & checksums
│
├── core/                        # Core data structures & logic
│   ├── object.go                # Obj struct (type, encoding, value, LRU clock)
│   ├── typeencoding.go          # Type/encoding helpers & assertions
│   ├── type_string.go           # String encoding deduction (INT / EMBSTR / RAW)
│   ├── store.go                 # In-memory key-value store + expiry map
│   ├── cmd.go                   # RedisCmd / RedisCmds types
│   ├── resp.go                  # RESP encoder & decoder
│   ├── eval.go                  # Command evaluation & dispatch
│   ├── expire.go                # Active expiry via sampling
│   ├── eviction.go              # Eviction strategies (first, random, LRU)
│   ├── evictionPool.go          # LRU eviction pool (priority queue)
│   ├── aof.go                   # AOF dump logic
│   ├── events.go                # Shutdown handler
│   ├── comm.go                  # FD-based ReadWriter for epoll
│   └── stats.go                 # Keyspace statistics
│
├── server/                      # Server implementations
│   ├── async_tcp.go             # Epoll-based async TCP server (default)
│   └── sync_tcp.go              # Blocking TCP server (net package)
│
├── internals/
│   └── config/
│       └── config.go            # Global configuration & CLI flags
│
└── bin/
    └── appendonly.aof           # Example AOF persistence file
```

---

## Getting Started

### Docker (recommended)

Pull and run the pre-built image from [Docker Hub](https://hub.docker.com/r/ishanksharma16/go-redis):

```bash
docker pull ishanksharma16/go-redis:latest
docker run -p 7379:7379 ishanksharma16/go-redis
```

### Prerequisites

- Go **1.26.3** or later (if building from source)

### Install & Run

```bash
# Clone the repository
git clone https://github.com/IshankSharma2178/go-redis.git
cd go-redis

# Build the binary
go build -o go-redis .

# Start the server (default port 7379)
./go-redis
```

Or run directly without building:

```bash
go run main.go
```

### Connect with redis-cli

```bash
redis-cli -p 7379
```

```redis
127.0.0.1:7379> PING
PONG
127.0.0.1:7379> SET mykey hello EX 10
OK
127.0.0.1:7379> GET mykey
"hello"
127.0.0.1:7379> TTL mykey
(integer) 7
127.0.0.1:7379> INCR counter
(integer) 1
127.0.0.1:7379> INCR counter
(integer) 2
127.0.0.1:7379> BGREWRITEAOF
OK
```

### Command-Line Flags

| Flag         | Default                   | Description                |
| ------------ | ------------------------- | -------------------------- |
| `-host`      | `0.0.0.0`                 | Host to bind               |
| `-port`      | `7379`                    | Port to listen on          |
| `-aof-file`  | `./bin/appendonly.aof`    | Path to AOF persistence    |

Example:

```bash
go-redis -port 6380 -aof-file /tmp/redis.aof
```

### Using netcat (RESP raw)

```bash
echo -e '*1\r\n$4\r\nPING\r\n' | nc localhost 7379
```

---

## How It Works

### Server

The default server uses **Linux epoll** for async, non-blocking TCP I/O. A single goroutine manages the event loop, accepting connections and reading commands in a state machine (`WAITING` → `BUSY` → `SHUTTING_DOWN`). A cron goroutine runs every second to sample and delete expired keys.

### Storage

Keys are stored in a Go `map[string]*Obj`. Each `Obj` holds a `TypeEncoding` byte (type in upper nibble, encoding in lower) and the actual value. Expiry timestamps are tracked in a separate `map[*Obj]uint64` and checked on every read access plus background sampling.

### Eviction

Three strategies are available:
- **simple-first** — iterates keys and deletes the first found
- **allkeys-random** — picks random keys until the eviction ratio is met
- **allkeys-lru** — maintains a pool of candidate keys sorted by idle time, evicting the least-recently-used

### Persistence

On `BGREWRITEAOF` or graceful shutdown, all keys are serialized in RESP format and written to the configured AOF file, providing crash recovery on restart.

---

## License

This project is open source and available under the [MIT License](LICENSE).
