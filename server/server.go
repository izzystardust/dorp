package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"text/template"

	"github.com/millere/dorp"
	"golang.org/x/crypto/nacl/secretbox"
)

type states struct {
	Door  dorp.State
	Light dorp.State
	sync.Mutex
}

func (s *states) Set(door, light dorp.State) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Door = door
	s.Light = light
}

var CurrentState = states{
	Door:  dorp.Negative,
	Light: dorp.Negative,
}

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
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handler)

	go ListenClients(conf.StatusPort, &key)

	http.ListenAndServe(fmt.Sprintf(":%d", conf.WebPort), nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index.html").ParseFiles("html/index.html")
	if err != nil {
		log.Println("Error executing template, using fallback output")
		fmt.Fprintf(w, "Door state: %s\nLight state: %s", CurrentState.Door, CurrentState.Light)
	}
	CurrentState.Lock()
	defer CurrentState.Unlock()
	err = t.Execute(w, CurrentState)
	if err != nil {
		log.Printf("Error executing template: %s", err)
	}
}

func ListenClients(port uint16, key *[32]byte) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		log.Println("Client connected:", conn.RemoteAddr())
		HandleClient(conn, key)
	}
}

func HandleClient(client net.Conn, key *[32]byte) {
	for {
		message, nonce := CreateNextReply(key)
		client.Write(message)
		var rawReply [256]byte
		n, err := client.Read(rawReply[:])
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected:", client.RemoteAddr())
				return
			}
			log.Println("bad:", err)
			client.Close()
			return
		}
		var reply []byte
		reply, ok := secretbox.Open(reply, rawReply[:n], nonce, key)
		if !ok {
			log.Println("not okay, dude")
			client.Close()
			return
		}
		if err := ProcessReply(reply); err != nil {
			log.Println(err)
			client.Close()
			return
		}
	}
}

func ProcessReply(reply []byte) error {
	if len(reply) != 2 {
		return dorp.ErrWrongNumberOfStates
	}
	doorState := dorp.State(reply[0])
	lightState := dorp.State(reply[1])
	CurrentState.Set(doorState, lightState)
	log.Println("Setting door:", doorState, "light:", lightState)
	return nil
}

func CreateNextReply(key *[32]byte) ([]byte, *[24]byte) {
	nonce, err := dorp.GenerateNonce(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	nextNonce, err := dorp.GenerateNonce(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	var box []byte
	box = secretbox.Seal(box, nextNonce[:], &nonce, key)
	return append(box, nonce[:]...), &nextNonce
}
