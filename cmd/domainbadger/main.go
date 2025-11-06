package main

import (
	"context"
	"fmt"
	"os"

	"github.com/emmanuelay/badger/internal/app"
	"github.com/emmanuelay/badger/internal/config"
)

// These will be set at build time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	ctx := context.Background()

	cfg, err := config.GetConfigurationFromArguments(version, commit, date)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	app.Run(ctx, cfg)
}
