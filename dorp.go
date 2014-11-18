package dorp

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/nacl/secretbox"
)

// A SetMessage is the Go representation of the JSON message
// sent to set states
type SetMessage struct {
	DoorState  string
	LightState string
}

// A State is  a binary condition of the door or lights.
type State int

const (
	Open State = iota
	Closed
	On
	Off
)

// PADDING_SIZE is the amount of padding used before the auth
// token in the message. The padding is used so that the same
// token encrypted with the same key produces different results
const PADDING_SIZE = 6

const DELIMITER = "/"

// String implements Stringer on States.
func (s State) String() string {
	switch s {
	case Open:
		return "Open"
	case Closed:
		return "Closed"
	case On:
		return "On"
	case Off:
		return "Off"
	default:
		panic("BAD STATE")
	}
}

// Encrypt takes a message and converts it to a base64 encoding
// of the encrypted string, followed by a separator, followed
// by the nonce
func Encrypt(key [32]byte, text []byte) (string, error) {
	var box []byte
	var nonce [24]byte
	n, err := rand.Reader.Read(nonce[:])
	if n != 24 {
		return "", fmt.Errorf("encrypt: unable to read 24 random bytes for nonce")
	}
	if err != nil {
		return "", err
	}
	box = secretbox.Seal(box[:0], text, &nonce, &key)
	return strings.Join([]string{
		base64.StdEncoding.EncodeToString(box),
		base64.StdEncoding.EncodeToString(nonce[:]),
	}, DELIMITER), nil
}

// Decrypt takes the base64 encoded box and nonce separated by DELIMITER
// and returns the opened box or an error
func Decrypt(data string, key [32]byte) ([]byte, error) {
	var nonce [24]byte
	var opened []byte
	parts := strings.Split(data, DELIMITER)
	if len(parts) != 2 {
		return nil, fmt.Errorf("decrypt: data contains too many delimiters")
	}
	box, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("decrypt box: %s", err)
	}
	nonceS, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decrypt nonce: %s", err)
	}
	if n := copy(nonce[:], nonceS); n != 24 {
		return nil, fmt.Errorf("decrypt: nonce has incorrect length")
	}
	opened, ok := secretbox.Open(opened, box, &nonce, &key)
	if !ok {
		return nil, fmt.Errorf("decrypt: failed to open box")
	}
	return opened, nil
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
