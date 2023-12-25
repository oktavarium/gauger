package main

import (
	"fmt"

	"github.com/oktavarium/go-gauger/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("error on running server: %w", err))
	}
}
