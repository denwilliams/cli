package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/denwilliams/cli/internal/aussupermarkets/cli"
	"github.com/denwilliams/cli/internal/aussupermarkets/providers"
)

func main() { os.Exit(run()) }
func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return cli.Run(ctx, os.Args[1:], os.Stdout, os.Stderr, providers.All())
}
