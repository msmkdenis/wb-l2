package main

import (
	"os"
	"os/signal"
	"syscall"

	"dev11/internal/app/calendar"
)

func main() {
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	calendar.Run(quitSignal)
}
