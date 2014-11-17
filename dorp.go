package dorp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// A SetMessage is the Go representation of the JSON message
// sent to set states
type SetMessage struct {
	DoorState  string
	LightState string
	Auth       string
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

// Encrypt encrypts text with key using AES and
// returns the encrypted text.
func Encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

// Decrypt decrypts text with key using AES and
// returns the decrypted text
func Decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// EncodeAuthToken encodes authentication token token with the given key
// and returns the base64 encoded token
func EncodeAuthToken(key, token string) (string, error) {
	padding := make([]byte, PADDING_SIZE)
	if _, err := rand.Read(padding); err != nil {
		return "", err
	}

	rawAuth, err := Encrypt(
		[]byte(key),
		append(padding, []byte(token)...),
	)
	if err != nil {
		return "", err
	}

	auth := base64.StdEncoding.EncodeToString(rawAuth)
	return string(auth), nil
}

// AuthIsValid compares expected to the value of got decrypted with key,
// removing the padding bytes from the front of got
func AuthIsValid(key, got, expected string) (bool, error) {
	auth, err := base64.StdEncoding.DecodeString(got)
	if err != nil {
		return false, fmt.Errorf("Couldn't decode base64: %s", err)
	}
	plain, err := Decrypt([]byte(key), []byte(auth))
	if err != nil {
		return false, fmt.Errorf("Couldn't decrypt token (bad token): %s", err)
	}

	// remove PADDING_SIZE bytes from the front of decrypted plain and compare
	if string(plain)[PADDING_SIZE:] == expected {
		return true, nil
	}

	return false, fmt.Errorf("Mismatched authentication tokens")
}
