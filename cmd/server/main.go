package main

import (
	"fmt"

	"github.com/oktavarium/go-gauger/internal/server"
)

var buildVersion string
var buildDate string
var buildCommit string

const NA string = "N/A"

func main() {
	printBuildInfo()
	if err := server.Run(); err != nil {
		panic(fmt.Errorf("error on running server: %w", err))
	}
}

func printBuildInfo() {
	if len(buildVersion) == 0 {
		buildVersion = NA
	}

	if len(buildDate) == 0 {
		buildDate = NA
	}

	if len(buildCommit) == 0 {
		buildCommit = NA
	}

	fmt.Printf(
		"Build version: %s\n Build data: %s\n Build commit: %s\n",
		buildVersion,
		buildDate,
		buildCommit,
	)
}
