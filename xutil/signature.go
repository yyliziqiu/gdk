package xutil

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	ErrTokenFormatError = errors.New("token format error")
	ErrTokenExpired     = errors.New("token expired")
	ErrInvalidSignature = errors.New("invalid signature")
)

func GenerateSignature(salt string, key string) string {
	return fmt.Sprintf("%s%x", salt, sha256.Sum256([]byte(salt+key)))
}

func VerifySignature(token string, key string, saltLength int) error {
	if len(token) != saltLength+64 {
		return ErrTokenFormatError
	}

	salt := token[:saltLength]

	if token != GenerateSignature(salt, key) {
		return ErrInvalidSignature
	}

	return nil
}

func GenerateTimeSignature(key string) string {
	salt := strconv.FormatInt(time.Now().Unix(), 10)

	return GenerateSignature(salt, key)
}

func VerifyTimeSignature(token string, ttl time.Duration, key string) error {
	if len(token) != 74 {
		return ErrTokenFormatError
	}

	saltLength := 10

	timestamp, err := strconv.Atoi(token[:saltLength])
	if err != nil {
		return ErrTokenFormatError
	}

	if time.Now().Add(-ttl).Unix() > int64(timestamp) {
		return ErrTokenExpired
	}

	return VerifySignature(token, key, saltLength)
}
