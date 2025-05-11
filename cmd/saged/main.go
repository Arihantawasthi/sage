package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Arihantawasthi/sage.git/cmd/saged/handlers"
	"github.com/Arihantawasthi/sage.git/internal/config"
	"github.com/Arihantawasthi/sage.git/internal/logger"
	"github.com/Arihantawasthi/sage.git/internal/models"
	"github.com/Arihantawasthi/sage.git/internal/spmp"
)

func main() {
    logger, err := logger.NewSlogLogger(models.LogFilePath)
    if err != nil {
		fmt.Fprintf(os.Stderr, "error while creating a logger: %s\n", err)
		os.Exit(1)
    }
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading config file: %s\n", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cmdMux := spmp.NewCommandMux()
    handler := handlers.NewHandler(config, *logger)
	cmdMux.HandleCommand(spmp.TypeStart, handler.HandleStartService)
	cmdMux.HandleCommand(spmp.TypeList, handler.HandleListServices)

	spmpServer := spmp.NewSPMPServer(cmdMux)
	go func(ctx context.Context) {
		spmpServer.ListenAndServe(ctx)
	}(ctx)

	<-ctx.Done()
	fmt.Println("Exiting...")
}
