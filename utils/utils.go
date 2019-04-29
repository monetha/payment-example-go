package utils

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// CreateCtrlCContext create a context which listens to interrupt command
func CreateCtrlCContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
		<-sigChan
		log.Println("got interrupt signal")
	}()

	return ctx
}
