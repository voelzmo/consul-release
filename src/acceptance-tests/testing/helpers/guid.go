package helpers

import (
	"crypto/rand"
	"fmt"
)

func NewGUID() (string, error) {
	var buf [16]byte

	_, err := rand.Read(buf[:])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:]), nil
}
