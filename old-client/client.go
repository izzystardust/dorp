package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/millere/dorp"
)

type Config struct {
	Server string
	Key    string
	Token  string
}

func main() {
	confFile := flag.String("f", "dorp.toml", "Configuration file")
	newDState := flag.String("d", "", "Door state to set")
	newLState := flag.String("l", "", "Light state to set")
	flag.Parse()
	if !(strings.ToLower(*newDState) == "open" || strings.ToLower(*newDState) == "closed") {
		fmt.Println("Door states must be open or closed. Set using -d flag")
		fmt.Println("Given d: ", *newDState)
		return
	}
	if !(strings.ToLower(*newLState) == "on" || strings.ToLower(*newLState) == "off") {
		fmt.Println("Light states must be on or off. Set using -l flag")
		fmt.Println("Given l: ", *newLState)
	}

	var c Config
	_, err := toml.DecodeFile(*confFile, &c)
	if err != nil {
		log.Fatal(err)
	}

	key := []byte(c.Key)
	auth, err := dorp.Encrypt(key, []byte(c.Token))
	authKey := base64.StdEncoding.EncodeToString(auth)
	if err != nil {
		log.Fatal(err)
	}

	var message bytes.Buffer
	encoder := json.NewEncoder(&message)
	encoder.Encode(dorp.SetMessage{
		DoorState:  *newDState,
		LightState: *newLState,
		Auth:       string(authKey),
	})
	_, err = http.Post(c.Server+"/set", "text/json", &message)
	if err != nil {
		log.Fatal(err)
	}
}
