package main

import (
	"go-image-compression/internal/app"
)

func main() {
	if err := app.MustRun(); err != nil {
		panic(err)
	}
}
