package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/millere/dorp"
)

var currentDoorState = dorp.Closed
var currentLightState = dorp.Off

var c Config

func main() {
	conffile := flag.String("f", "dorp.toml", "Configuration file")
	flag.Parse()
	conf, err := ReadConfig(*conffile)
	if err != nil {
		log.Fatal(err)
	}
	c = conf
	http.HandleFunc("/", handler)
	http.HandleFunc("/set", setState)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Door state: %s\nLight state: %s", currentDoorState, currentLightState)
}

func setState(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var m dorp.SetMessage
	err := decoder.Decode(&m)
	if err != nil {
		log.Println("Bad data given")
		return
	}

	authValid, err := dorp.AuthIsValid(c.Key, m.Auth, c.Token)
	if err != nil {
		log.Println(err)
		return
	}
	if !authValid {
		log.Println("Invalid auth without error. A bug exists")
		return
	}

	if strings.ToLower(m.DoorState) == "open" {
		currentDoorState = dorp.Open
	} else if strings.ToLower(m.DoorState) == "closed" {
		currentDoorState = dorp.Closed
	} else {
		log.Println("Bad state recvd: ", m.DoorState)
		return
	}
	if strings.ToLower(m.LightState) == "on" {
		currentLightState = dorp.On
	} else if strings.ToLower(m.LightState) == "off" {
		currentLightState = dorp.Off
	} else {
		log.Println("Bad state recvd: ", m.LightState)
		return
	}
	log.Println("Set states:", currentDoorState, currentLightState)
}
