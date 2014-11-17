# Dorp Protocol

The server listens on /set for JSON messages of the following format:

```JSON
{
	"DoorState": "open|closed",
	"LightState": "on|off",
	"Auth": "auth token",
}
```

DoorState must be either "open" or "closed" and sets the door state to that.

LightState must be either "on" or "off" and sets the light state to that.

Auth TODO
