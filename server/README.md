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

## Building and Deploying

Building is trivial: from the server directory, run `go build`.

To deploy, place the compiled `server` binary as well as the `html/` and `static/` directories on the server.
Run the server program to start serving. Recommendation: Do it in a tmux, so you can disconect.
