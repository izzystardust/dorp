package dorp

import (
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

// A SetMessage is the Go representation of the JSON message
// sent to set states
type SetMessage struct {
	DoorState  string
	LightState string
}

// A State is  a binary condition of the door or lights.
type State byte

//go:generate stringer -type=State

// Positive and Negative are the two possible states
const (
	Negative State = iota
	Positive
)

// ErrWrongNumberOfStates may be returned if a function takes a slice of states
// and got the wrong number. This probably will only occur deserializing network data.
var ErrWrongNumberOfStates = errors.New("incorrect state count")

// GenerateNonce creats a 24 byte nonce from the source of randomness rand
func GenerateNonce(rand io.Reader) ([24]byte, error) {
	var nonce [24]byte
	var n int
	for i := 0; i < 3; i++ {
		n, _ = rand.Read(nonce[:])
		if n == 24 {
			break
		}
	}
	if n != 24 {
		return nonce, fmt.Errorf("encrypt: unable to read 24 random bytes for nonce (read %d)", n)
	}
	return nonce, nil
}

// KeyToByteArray converts a key to the [32]bytes required by nacl.secretbox
func KeyToByteArray(key string) ([32]byte, error) {
	var k [32]byte
	if len(key) != 32 {
		return k, fmt.Errorf("Key must be 32 bytes (characters) long")
	}
	n := copy(k[:], []byte(key))
	if n != 32 {
		return k, fmt.Errorf("Copying key failed")
	}
	return k, nil
}

// ProcessNonceMessage takes the message from the server and the shared key
// and returns the next nonce the server expects
func ProcessNonceMessage(message *[64]byte, key *[32]byte) (*[24]byte, error) {
	var nonce [24]byte
	copy(nonce[:], message[64-24:])
	var nextNonce []byte
	var ok bool
	nextNonce, ok = secretbox.Open(nextNonce, message[:64-24], &nonce, key)
	if !ok {
		return nil, fmt.Errorf("Unable to open box")
	}
	n := copy(nonce[:], nextNonce)
	if n != 24 {
		return nil, fmt.Errorf("Recvd nonce has incorrect length")
	}
	return &nonce, nil
}

// Thought of a haiku
// I may as well leave it here
// To be found by you
