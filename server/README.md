# The Dorp Server

The dorp server does two things:

1. Serve current status on `/`
2. Listen for encrypted update messages the configured tcp port

## Configuration
Dorp is configured using a toml file. It should look something like
```toml
key = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
statusport = 13699
webport = 8080
```

`key` should be a pre-shared key with the client. It must be 32 bytes long.

`statusport` should be the port to listen on for status updates.

`webport` is the port that the server will serve pages on.

For documentation on the message format, see `protocol.md`.

## Building
The Dorp server uses [go-bindata](https://github.com/jteeuwen/go-bindata) to include static files in the distributed binary.
If no changes to the static files in `html/` and `static/` are needed, use the go tool as normal.
If the static files are modified, first run `go-bindata html/ static/` to regenerate the static resources.
The `-debug` flag causes `go-bindata` to produce a stub that reads from the filesystem for rapid development.
