package main

import (
	"errors"
	"flag"
)

var flagRunAddr string

func parseFlags() error {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080",
		"address and port of server in notaion address:port")
	flag.Parse()

	if len(flag.Args()) > 0 {
		return errors.New("unrecognised flags")
	}
	return nil
}
