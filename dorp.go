package dorp

import (
	"errors"
	"fmt"
	"io"
)

// A SetMessage is the Go representation of the JSON message
// sent to set states
type SetMessage struct {
	DoorState  string
	LightState string
}

// A State is  a binary condition of the door or lights.
type State byte

const (
	Positive State = iota
	Negative
)

// String implements Stringer on States.
func (s State) String() string {
	switch s {
	case Positive:
		return "✔"
	case Negative:
		return "✘"
	default:
		panic("BAD STATE")
	}
}

var IncorrectStateCount = errors.New("incorrect state count")

// GenerateNonce creats a 24 byte nonce from the source of randomness rand
func GenerateNonce(rand io.Reader) ([24]byte, error) {
	var nonce [24]byte
	n, err := rand.Read(nonce[:])
	if n != 24 {
		return nonce, fmt.Errorf("encrypt: unable to read 24 random bytes for nonce", n, nonce[:])
	}
	if err != nil {
		return nonce, err
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
