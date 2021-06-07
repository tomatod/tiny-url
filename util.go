package main

import (
	"crypto/rand"
	"errors"
)

const lettersForRandomStr = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func MakeRandomStr(digit uint32) (string, error) {
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("unexpected error...")
	}

	var result string
	for _, v := range b {
		result += string(lettersForRandomStr[int(v)%len(lettersForRandomStr)])
	}
	return result, nil
}
