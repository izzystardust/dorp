# Dorp Protocol

The server listens on /set for JSON messages of the following format:

```JSON
{
	"DoorState": "open|closed",
	"LightState": "on|off",
	"Auth": "auth token",
}
```

## DoorState
DoorState must be either "open" or "closed" and sets the door state to that.

## LightState
LightState must be either "on" or "off" and sets the light state to that.

## Auth
The authentication token is more complex. The server and client must have two shared secrets:
an AES key and a shared token.

Auth is currently using a broken protocol, as the same auth token sent twice will work.
TODO: Switch to nacl secretbox, golang.org/x/crypto/nacl/secretbox
