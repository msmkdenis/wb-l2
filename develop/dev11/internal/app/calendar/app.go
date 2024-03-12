package calendar

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"dev11/internal/config"
	"dev11/internal/delivery/http/handlers"
	"dev11/internal/delivery/http/middleware"
	"dev11/internal/repository"
	"dev11/internal/service"
)

func Run(quitSignal chan os.Signal) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	cfg := config.New()

	mux := http.NewServeMux()

	calendarRepository := repository.NewCalendar()
	calendarService := service.NewCalendar(calendarRepository)

	requestLogger := middleware.NewRequestLogger()
	handlers.NewCalendar(calendarService, mux, requestLogger)

	httpServ := &http.Server{
		Addr:    cfg.URLServer,
		Handler: mux,
	}

	httpServerCtx, httpServerStopCtx := context.WithCancel(context.Background())

	go func() {
		slog.Info("staring server", slog.String("addr", cfg.URLServer))
		err := httpServ.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start server", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan struct{})
	go func() {
		<-quitSignal
		close(quit)
	}()

	go func() {
		<-quit

		// Shutdown signal with grace period of 10 seconds
		shutdownCtx, cancel := context.WithTimeout(httpServerCtx, 10*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				slog.Error("graceful shutdown timed out.. forcing exit.")
				os.Exit(1)
			}
		}()

		// Trigger graceful shutdown
		slog.Info("Initiating graceful shutdown")
		if err := httpServ.Shutdown(shutdownCtx); err != nil {
			slog.Error("Failed to shutdown server", slog.String("err", err.Error()))
			os.Exit(1)
		}
		httpServerStopCtx()
	}()

	<-httpServerCtx.Done()
}
