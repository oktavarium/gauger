package main

import (
	"fmt"

	"github.com/oktavarium/go-gauger/internal/agent"
)

func main() {
	if err := agent.Run(); err != nil {
		panic(fmt.Errorf("error on running agent: %w", err))
	}
}
