package main

import (
	"os"
	"github.com/slimming/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
