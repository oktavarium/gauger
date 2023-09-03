package main

import (
	"errors"
	"flag"
)

var flagEndpointAddr string
var flagReportInterval int
var flagPollInterval int

func parseFlags() error {
	flag.StringVar(&flagEndpointAddr, "a", "localhost:8080",
		"address and port of server's endpoint in notaion address:port")
	flag.IntVar(&flagReportInterval, "r", 10,
		"report interval in seconds")
	flag.IntVar(&flagPollInterval, "p", 2,
		"apoll interval in seconds")
	flag.Parse()

	if len(flag.Args()) > 0 {
		return errors.New("unrecognised flags")
	}
	return nil
}
