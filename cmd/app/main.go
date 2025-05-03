package main

import (
	"log"

	"go-image-compression/internal/app"
)

func main() {
	err := app.MustRun()
	if err != nil {
		log.Fatalf("cmd/app/main.go/: %v", err)
	}
}
