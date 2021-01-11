package main

import (
	"fmt"
	"os"

	"github.com/emmanuelay/domainsearch/internal/app"
	"github.com/emmanuelay/domainsearch/internal/config"
)

func main() {
	// Get config from commandline arguments
	cfg, err := config.GetConfigurationFromArguments()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	app.Run(cfg)
}
