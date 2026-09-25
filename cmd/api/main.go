package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/fykzi/url-shortener/internal/config"
	"github.com/fykzi/url-shortener/internal/db"
	"github.com/fykzi/url-shortener/internal/handler"
	"github.com/fykzi/url-shortener/internal/logger"
	"github.com/fykzi/url-shortener/internal/repository/sqlite"
	"github.com/fykzi/url-shortener/internal/service"
)

func main() {
	slog.Info("Load config")
	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Setup logger")
	log, err := logger.SetupLogger(cfg.Log.Level)
	if err != nil {
		slog.Error("setup logger failed", "error", err)
		os.Exit(1)
	}

	if err := run(cfg, log); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run(cfg *config.Config, baselog *slog.Logger) error {
	database, err := db.New(cfg.DB.Path)
	if err != nil {
		return fmt.Errorf("database connection: %w", err)
	}
	defer database.Close()
	slog.Info("database connected", "path", cfg.DB.Path)

	// repository
	repo := sqlite.NewLinkRepository(database, baselog.With("component", "repo"))

	// service
	svc := service.NewLinkService(repo, cfg.Server.BaseURL(), baselog.With("component", "service"))

	// handlers
	linkH := handler.NewLinkHandler(svc, baselog.With("component", "link_handler"))
	healthH := handler.NewHealthHandler(baselog.With("component", "health_handler"))

	// router
	mux := handler.RegisterRoutes(linkH, healthH)

	baselog.Info("Server starting", "addr", cfg.Server.Addr())
	if err := http.ListenAndServe(cfg.Server.Addr(), mux); err != nil {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}
