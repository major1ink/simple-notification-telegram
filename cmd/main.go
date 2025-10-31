package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/major1ink/simple-notification-telegram/internal/app"
)

//go:embed VERSION
var version []byte

//go:generate go run ../script/version/

func main() {
	args := os.Args
	for _, arg := range args {
		if arg == "--version" {
			fmt.Println(string(version))
			os.Exit(0)
		}
	}

	var configPath string
	flag.StringVar(&configPath, "configPath", "", "path to config file")

	flag.Parse()

	if configPath != "" {
		if err := os.Setenv("CONFIG_PATH", configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to set CONFIG_PATH: %v\n", err)
			os.Exit(1)
		}
	}

	ctx := context.Background()
	a, err := app.New(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create app: %v\n", err)
		os.Exit(1)
	}

	if err := a.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "App failed: %v\n", err)
		os.Exit(1)
	}
}
