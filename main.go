package main

import (
	"os"

	"nexus-sds.com/netcheck/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
