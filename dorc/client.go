package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/millere/dorp"
	"github.com/stianeikeland/go-rpio"
	"golang.org/x/crypto/nacl/secretbox"
)

type Config struct {
	Server   string
	Port     uint16
	Key      string
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
	server, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Server, c.Port))
	if err != nil {
		log.Fatal(err)
	}

	lVal, l := MonitorPin(lp, 5*time.Second)
	dVal, d := MonitorPin(dp, 5*time.Second)

	ticker := time.Tick(5 * time.Minute)
	for {
		err := SendUpdate(dVal, lVal, server, &key)
		if err == nil {
		} else {
			log.Println("Error sending new states: ", err)
		}
		select {
		case newL := <-l:
			lVal = newL
		case newD := <-d:
			dVal = newD
		case <-ticker:
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
func SendUpdate(door, light rpio.State, server net.Conn, key *[32]byte) error {
	var message [64]byte
	_, err := server.Read(message[:])
	if err != nil {
		return err
	}
	nonce, err := dorp.ProcessNonceMessage(&message, key)
	if err != nil {
		return err
	}
	update := []byte{byte(door), byte(light)}
	var cipher []byte
	cipher = secretbox.Seal(cipher, update, nonce, key)
	n, err := server.Write(cipher)
	if n != len(cipher) {
		panic("This shouldn't be happening. This can't happen")
	}
	return err
}
