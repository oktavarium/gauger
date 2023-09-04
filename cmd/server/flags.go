package main

import (
	"errors"
	"flag"
	"os"
)

var flagRunAddr string

const (
	envRunAddrName string = "ADDRESS"
)

func parseFlags() error {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080",
		"address and port of server in notaion address:port")
	flag.Parse()

	if len(flag.Args()) > 0 {
		return errors.New("unrecognised flags")
	}

	envRunAddr := os.Getenv(envRunAddrName)
	if len(envRunAddr) != 0 {
		flagRunAddr = envRunAddr
	}

	return nil
}
