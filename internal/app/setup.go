package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rmarasigan/warehouse-inventory-management/api"
	"github.com/rmarasigan/warehouse-inventory-management/internal/app/config"
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/mysql"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/log"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/trail"
)

const (
	FGNormal   = "\x1b[0m"
	FGRedB     = "\x1b[31;1m"
	FGMagentaB = "\x1b[35;1m"
)

const message = "" +
	FGMagentaB + `              ` + FGRedB + `    ______   _____` + "\n" + FGNormal +
	FGMagentaB + `   ___  ___  ___` + FGRedB + ` /  /    \/     \` + "\n" + FGNormal +
	FGMagentaB + `  /  / /  / /  /` + FGRedB + `/  /  __   __   /` + "\n" + FGNormal +
	FGMagentaB + ` /  / /  / /  /` + FGRedB + `/  /  / /  / /  /` + "\n" + FGNormal +
	FGMagentaB + `/  /_/  /_/  /` + FGRedB + `/  /  / /  / /  /` + "\n" + FGNormal +
	FGMagentaB + `\_____,_____/` + FGRedB + `/__/__/ /__/ /__/` + FGNormal + "\n\n"

func StartServer() {
	log.Init()
	defer log.Panic()

	// Load the application configuration
	cfg, err := config.Load("wim-config.yaml")
	if err != nil {
		panic(err)
	}

	// Connect to the MySQL database
	mysql.Connect()

	// Set-up the HTTP server
	serverAddress := cfg.ServerAddress()
	server := &http.Server{
		Addr:    serverAddress,
		Handler: http.DefaultServeMux,
	}

	// Register the applications handler
	http.HandleFunc("/api/", api.Handler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Handle a graceful shutdown
	go func() {
		// Wait for interrupt signal to gracefully shutdown
		<-quit

		trail.Info("Shutting down server...")
		mysql.Close()

		// Create a context with a timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Attempt to shutdown the server gracefully
		err := server.Shutdown(ctx)
		if err != nil {
			log.Error(err, "failed to shutdown the server")
		}
	}()

	// Start the server
	trail.Info("Initializing %s Server at %s\n%s", cfg.AppName(), serverAddress, message)

	err = server.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
