# go-http-server

A minimal HTTP server built from scratch in Go, using raw TCP sockets without the standard `net/http` package. Designed for learning and experimentation with the HTTP protocol and low‑level network programming.

## Features

- **Custom HTTP parsing**: Manual request parsing and response serialization.
- **Persistent connections**: Supports HTTP/1.1 keep‑alive by default; respects `Connection: close` header.
- **Gzip compression**: Honors `Accept-Encoding: gzip` and compresses responses when requested.
- **Echo endpoint**: `/echo/{message}` returns the text after `/echo/` as the response body.
- **User-Agent reflector**: `/user-agent` returns the client's `User-Agent` header.
- **Static file serving & upload**:

  - `GET /files/{path}`: Serves the contents of the given file path.
  - `POST /files/{path}`: Writes the request body to the given file path.

- **404 handling**: Returns a 404 Not Found for unknown paths.

## Prerequisites

- Go 1.18 or later installed.
- A POSIX‑compatible shell (for Makefile targets).

## Getting Started

Clone the repository:

```bash
git clone https://github.com/DylanPIna/go-http.git
cd go-http
```

### Build

To compile the server:

```bash
make build
```

This produces the binary `go-http-server` (or `go-http-server.exe` on Windows).

Or, build directly with `go`:

```bash
go build -o go-http-server cmd/server/main.go
```

### Run

Start the server:

```bash
./go-http-server
```

By default, it listens on TCP port `4221` on all interfaces (`0.0.0.0:4221`).

### Usage

#### Root `/`

```bash
curl http://localhost:4221/
```

Returns an HTTP 200 OK with an empty body.

#### Echo `/echo/{message}`

```bash
curl http://localhost:4221/echo/hello
```

Response body: `hello`

#### User-Agent `/user-agent`

```bash
curl -A "MyClient/1.0" http://localhost:4221/user-agent
```

Response body: `MyClient/1.0`

#### Static File Serving `/files/{path}`

```bash
# Download a file:
curl http://localhost:4221/files/path/to/file.txt

# Upload a file:
curl -X POST --data-binary @localfile.txt http://localhost:4221/files/uploaded.txt
```

#### Gzip Compression

Requests with `Accept-Encoding: gzip` will receive compressed responses:

```bash
curl -H "Accept-Encoding: gzip" http://localhost:4221/echo/test --compressed
```

### Testing

Run the server’s unit tests:

```bash
make test
```

Or directly:

```bash
go test -v ./cmd/server
```

## Project Structure

```
.
├── cmd               # Entry point(s)
│   └── server        # Server implementation and tests
│       ├── main.go
│       └── main_test.go
├── internal
│   ├── http        # HTTP parsing and response types
│   ├── utils       # Utility functions (e.g., GzipCompress)
│   └── client      # Client code used for testing
├── Makefile
└── README.md
```

## Configuration

- **Port**: Hardcoded to `4221`. Change in `cmd/server/main.go` if needed.
- **Logging**: Simple `fmt.Println` for request/response events.
