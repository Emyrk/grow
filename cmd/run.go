package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		os.Exit(1)
	}()

	_ = RootCmd.ExecuteContext(ctx)
}
