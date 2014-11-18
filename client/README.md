# Dorp Client

## Configuration
Dorp is configured using a toml file. It should look something like
```toml
server = "k2cc.clarkson.edu:8080"
key = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
token = "expected"
doorPin = 17
lightPin = 16
```

`server` should be the server and port the dorp server is running on.

`key` and `token` should be the pre-shared keys. `key` must be 32 characters long. `token` may be
any length.

`doorPin` and `lightPin` should be set to the bcm2835 pin number, not the physical pin number.
For example, physical pin 19 is bcm2835 pin 10 and should be pin 10 in the configuration.
