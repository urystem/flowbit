package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	bootstrap "marketflow/internal/app"
	"marketflow/internal/config"
)

func main() {
	ctxBack := context.Background()
	cfg := config.Load()

	app, err := bootstrap.InitApp(ctxBack, cfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Run(); err != nil && err != http.ErrServerClosed {
			slog.Error("❌", " Server error:", err)
			quit <- syscall.SIGTERM
		}
	}()
	// time.Sleep(20 * time.Second)
	// close(quit)
	<-quit // Ждём сигнал
	slog.Info("📦 Shutting down server...")

	// Контекст с таймаутом на завершение
	ctx, cancel := context.WithTimeout(ctxBack, 20*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		slog.Error("❌", " Server forced to shutdown: %v", err)
	}
	slog.Info("✅ Server exited properly")
}
