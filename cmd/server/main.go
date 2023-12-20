package main

import (
	"fmt"
	"net/http"

	_ "net/http/pprof"

	"github.com/oktavarium/go-gauger/internal/server"
)

func main() {
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!1111")
	go func() {
		err := http.ListenAndServe(":8888", nil)
		if err != nil {
			panic(fmt.Errorf("pprf error: %w", err))
		}
	}()

	if err := server.Run(); err != nil {
		panic(fmt.Errorf("error on running server: %w", err))
	}
}
