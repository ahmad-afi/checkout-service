package utils

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2/log"
)

func TestIDGenerator(t *testing.T) {
	for i := 0; i < 10; i++ {
		id, err := IDGenerator()
		if err != nil {
			log.Error(err)
		}
		fmt.Println(id)
	}
}
