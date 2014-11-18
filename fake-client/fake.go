package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/millere/dorp"
	"github.com/stianeikeland/go-rpio"
)

func main() {
	server := "http://k2cc.clarkson.edu:8080"
	key, err := dorp.KeyToByteArray("abcdefghijklmnopqrstuvwxyzabcdef")
	if err != nil {
		panic(err)
	}

	c := time.Tick(7 * time.Second)
	door := rpio.Low
	light := rpio.High
	for {
		err := SendUpdate(door, light, server, key)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Sent d:%d l:%d", door, light)
		}
		light, door = door, light
		<-c
	}
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
	_, err = http.Post(server+"/set", "text/plain", strings.NewReader(data))
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
