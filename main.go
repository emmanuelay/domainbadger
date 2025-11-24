package main

import (
	"fmt"
	"os"

	"github.com/emmanuelay/badger/cmd"
)

// These will be set at build time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := cmd.Execute(version, commit, date); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
