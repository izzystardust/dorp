package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/millere/dorp"
	"golang.org/x/crypto/nacl/secretbox"
)

func main() {
	server := "localhost"
	port := 13699
	key, err := dorp.KeyToByteArray("abcdefghijklmnopqrstuvwxyzabcdef")
	if err != nil {
		panic(err)
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server, port))
	if err != nil {
		log.Fatal(err)
	}

	door := dorp.Positive
	light := dorp.Negative

	var message [64]byte
	var nonce [24]byte
	var nextNonce []byte
	var cipher []byte
	for {
		_, err = conn.Read(message[:])
		if err != nil {
			log.Fatal(err)
		}
		copy(nonce[:], message[64-24:])
		nextNonce, _ = secretbox.Open(nextNonce[:0], message[:64-24], &nonce, &key)
		copy(nonce[:], nextNonce)
		reply := []byte{byte(door), byte(light)}
		cipher = secretbox.Seal(cipher[:0], reply, &nonce, &key)
		log.Println("Sending door:", door, "light:", light)
		conn.Write(cipher)
		time.Sleep(3 * time.Second)
		door, light = light, door
	}
}
