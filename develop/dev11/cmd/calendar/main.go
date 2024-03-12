package main

import (
	"dev11/internal/app/calendar"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	calendar.Run(quitSignal)
}
