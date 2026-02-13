package main

import (
	"calculate_product_packs/internal/config"
	"calculate_product_packs/internal/repository"
	httphandler "calculate_product_packs/internal/transport/http"
	"calculate_product_packs/internal/usecases"
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	cfg := config.NewConfig()

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	repo := repository.NewMemoryPackSizeRepository(cfg.PackSizes)
	calculatePacksUseCase := usecases.NewCalculatePacksUseCase(repo)
	packSizesUseCase := usecases.NewPackSizesUseCase(repo)

	handler := httphandler.NewPackCalculatorHandler(calculatePacksUseCase, packSizesUseCase)
	router := httphandler.NewRouter(handler, tmpl)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("server starting", "port", cfg.Port)
		serverErrors <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	case <-shutdown:
		slog.Info("shutdown signal received")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("graceful shutdown failed", "error", err)
			srv.Close()
		}
	}

	slog.Info("server stopped")
}
