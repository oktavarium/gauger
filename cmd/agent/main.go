package main

import (
	"fmt"

	"net/http"
	_ "net/http/pprof"

	"github.com/oktavarium/go-gauger/internal/agent"
)

func main() {
	go func() {
		err := http.ListenAndServe(":8082", nil)
		if err != nil {
			panic(fmt.Errorf("pprf error: %w", err))
		}
	}()

	if err := agent.Run(); err != nil {
		panic(fmt.Errorf("error on running agent: %w", err))
	}
}
