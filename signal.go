package main

import (
	"os"
	"os/signal"
	"syscall"
)

var sigC = make(chan os.Signal, 1)

func init() {
	signal.Notify(sigC, syscall.SIGINT)
}
