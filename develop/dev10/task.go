package main

import (
	"flag"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port
go-telnet mysite.ru 8080
go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные из сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

type config struct {
	url     string
	timeout time.Duration
}

func (c *config) parseFlags() {
	var url string
	flag.StringVar(&url, "u", "tcpbin.com:4242", "url")

	var timeout time.Duration
	flag.DurationVar(&timeout, "t", time.Second*10, "timeout in seconds")

	flag.Parse()
	c.url = url
	c.timeout = timeout
}

func main() {
	cfg := config{}
	cfg.parseFlags()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	client := NewTCPClient(cfg.url, cfg.timeout, io.NopCloser(os.Stdin), os.Stdout)

	slog.Info("Connecting to", "server", cfg.url)
	if err := client.Connect(); err != nil {
		slog.Error("failed to connect", "error", err)
		os.Exit(1)
	}
	defer func() {
		err := client.Close()
		slog.Error("failed to close", "error", err)
		os.Exit(1)
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			if err := client.Send(); err != nil {
				slog.Error("failed to send", "error", err)
				os.Exit(1)
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			if err := client.Receive(); err != nil {
				slog.Error("failed to receive", "error", err)
				os.Exit(1)
				return
			}
		}
	}()

	wg.Wait()

}
