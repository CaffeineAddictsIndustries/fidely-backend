package main

import (
	"fmt"
	"log"
	"os"

	"fidely-backend/internal/auth"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: go run ./cmd/hash-password \"plain-password\"")
	}

	manager := auth.NewDefaultPasswordManager()
	hash, err := manager.Hash(os.Args[1])
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	fmt.Println(hash)
}
