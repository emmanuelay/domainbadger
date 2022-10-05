package main

import (
	"context"
	"fmt"
	"os"

	"github.com/emmanuelay/badger/internal/app"
	"github.com/emmanuelay/badger/internal/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.GetConfigurationFromArguments()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	app.Run(ctx, cfg)
}
