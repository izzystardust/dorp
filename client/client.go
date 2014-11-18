package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/millere/dorp"
	"github.com/stianeikeland/go-rpio"
)

type Config struct {
	Server   string
	Key      string
	Token    string
	DoorPin  uint8
	LightPin uint8
}

func main() {
	confFile := flag.String("f", "dorp.toml", "Configuration file")
	flag.Parse()

	var c Config
	if _, err := toml.DecodeFile(*confFile, &c); err != nil {
		log.Fatal("Error reading config: ", err)
	}

	key, err := dorp.KeyToByteArray(c.Key)
	if err != nil {
		log.Fatal("Error converting key: ", err)
	}

	lp, dp, err := Init(c.DoorPin, c.LightPin)
	if err != nil {
		log.Fatal("GPIO init error: ", err)
	}
	defer rpio.Close()

	lVal, l := MonitorPin(lp, 5*time.Second)
	dVal, d := MonitorPin(dp, 5*time.Second)
	for {
		select {
		case newL := <-l:
			lVal = newL
		case newD := <-d:
			dVal = newD
		}
		err := SendUpdate(dVal, lVal, c.Server, key)
		if err == nil {
			log.Printf(
				"Sent values [d: %s, l:%s]\n",
				DoorStateToString(dVal),
				LightStateToString(lVal),
			)
		} else {
			log.Println("Error sending new states: ", err)
		}

	}
}

// Init initializes the rpio memory mapped IO and sets up the door and light
// monitoring pins. It returns the door and light pins, and a possible error.
// If error is non-nil, an issue occured configuring memory mapped IO and
// reading from the returned pins is undefined behavior
func Init(dp, lp uint8) (rpio.Pin, rpio.Pin, error) {
	d := rpio.Pin(dp)
	l := rpio.Pin(lp)
	if err := rpio.Open(); err != nil {
		return d, l, err
	}
	d.Input()
	l.Input()
	return d, l, nil
}

// MonirtorPin checks pin a for a change in state once every interval.
// If a state change occurs, the new state is sent to the returned channel
func MonitorPin(a rpio.Pin, interval time.Duration) (rpio.State, <-chan rpio.State) {
	notify := make(chan rpio.State)
	ticker := time.Tick(interval)
	oldState := a.Read()
	go func() {
		for _ = range ticker {
			newState := a.Read()
			if newState != oldState {
				notify <- newState
			}
			oldState = newState
		}
	}()
	return oldState, notify
}

// SendUpdate sends the states of door and light to the given server,
// encrypting the authentication token token with key.
func SendUpdate(door, light rpio.State, server string, key [32]byte) error {
	var message bytes.Buffer
	encoder := json.NewEncoder(&message)
	encoder.Encode(dorp.SetMessage{
		DoorState:  DoorStateToString(door),
		LightState: LightStateToString(light),
	})
	data, err := dorp.Encrypt(key, message.Bytes())
	if err != nil {
		return err
	}
	_, err = http.Post("http://"+server+"/set", "text/plain", strings.NewReader(data))
	return err
}

func DoorStateToString(d rpio.State) string {
	switch d {
	case rpio.Low:
		return "Closed"
	case rpio.High:
		return "Open"
	default:
		panic("State can't exist")
	}
}
func LightStateToString(d rpio.State) string {
	switch d {
	case rpio.Low:
		return "Off"
	case rpio.High:
		return "On"
	default:
		panic("State can't exist")
	}
}
