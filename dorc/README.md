# Dorp Client

## Configuration
Dorp is configured using a toml file. It should look something like
```toml
server = "k2cc.clarkson.edu"
port = 13699
key = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
doorPin = 17
lightPin = 16
```

`server` should be the server and port the dorp server is running on.

`key` and should be a pre-shared key. `key` must be 32 bytes long. 

`doorPin` and `lightPin` should be set to the bcm2835 pin number, not the physical pin number.
For example, physical pin 19 is bcm2835 pin 10 and should be pin 10 in the configuration.

See [here](http://raspberrypi.znix.com/hipidocs/topic_gpiopins.htm) for more details on physical->bcm2835 number mapping.
