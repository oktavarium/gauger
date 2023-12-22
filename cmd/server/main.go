package main

import (
	"fmt"

	// _ "net/http/pprof"

	"github.com/oktavarium/go-gauger/internal/server"
)

func main() {
	// func() {
	// 	err := http.ListenAndServe(":8888", nil)
	// 	if err != nil {
	// 		panic(fmt.Errorf("pprf error: %w", err))
	// 	}go
	// }()

	if err := server.Run(); err != nil {
		panic(fmt.Errorf("error on running server: %w", err))
	}
}
