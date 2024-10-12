package utils

import (
	"log"
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

func IDGenerator() (res string, err error) {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	ulidres, err := ulid.New(ms, entropy)
	if err != nil {
		log.Println("failed to create ulid ", err)
		return
	}

	return ulidres.String(), err
}
