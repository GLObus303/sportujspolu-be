package utils

import (
	"log"

	nanoid "github.com/matoous/go-nanoid/v2"
)

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	length   = 12
)

func GenerateUUID() string {
	uuid, err := nanoid.Generate(alphabet, length)
	if err != nil {
		log.Println("(GenerateUUID) nanoid.Generate", err)
	}

	return uuid
}
