package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/millere/dorp"
)

var currentDoorState = dorp.Closed
var currentLightState = dorp.Off

var c Config
var key [32]byte

func main() {
	conffile := flag.String("f", "dorp.toml", "Configuration file")
	flag.Parse()
	conf, err := ReadConfig(*conffile)
	if err != nil {
		log.Fatal(err)
	}
	c = conf
	key, err = dorp.KeyToByteArray(c.Key)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", handler)
	http.HandleFunc("/set", setState)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Door state: %s\nLight state: %s", currentDoorState, currentLightState)
}

func setState(w http.ResponseWriter, r *http.Request) {
	rawMessage, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Couldn't ioutil.Readall body?")
		return
	}
	cipherMessage := string(rawMessage)
	plainMessage, err := dorp.Decrypt(cipherMessage, key)

	decoder := json.NewDecoder(strings.NewReader(string(plainMessage)))
	var m dorp.SetMessage
	err = decoder.Decode(&m)
	if err != nil {
		log.Println("Bad data given")
		return
	}

	if strings.ToLower(m.DoorState) == "open" {
		currentDoorState = dorp.Open
	} else if strings.ToLower(m.DoorState) == "closed" {
		currentDoorState = dorp.Closed
	} else {
		log.Println("Bad state recvd: ", m.DoorState)
	}
	if strings.ToLower(m.LightState) == "on" {
		currentLightState = dorp.On
	} else if strings.ToLower(m.LightState) == "off" {
		currentLightState = dorp.Off
	} else {
		log.Println("Bad state recvd: ", m.LightState)
	}
	log.Println("Set states:", currentDoorState, currentLightState)
}
