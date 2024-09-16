package main

import (
	"context"
	"errors"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/api"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/api/handler"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/config"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/service"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/storage/postgres"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/logger"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()

	sLog := logger.NewLogger(cfg.Env)

	connDB, err := postgres.NewStorage(context.Background(), cfg)
	if err != nil {
		sLog.Error("failed to init storage", logger.Error(err))
		os.Exit(1)
	}
	defer connDB.CloseDB()

	svc := service.NewService(cfg, sLog, connDB)

	hand := handler.NewHandler(cfg, sLog, svc)
	router := api.SetUpRoutes(hand, sLog)

	listenAddr := net.JoinHostPort(cfg.Host, cfg.Port)
	sLog.Info("starting server...", logger.String("address", listenAddr))

	srv := http.Server{
		Addr:         listenAddr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			sLog.Error("failed to start server", logger.Error(err))
		}
	}()

	sLog.Info("server started")

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		sLog.Error("failed to stop server", logger.Error(err))
		return
	}

	sLog.Info("server stopped")
}
