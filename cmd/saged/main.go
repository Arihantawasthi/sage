package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Arihantawasthi/sage.git/internal/config"
	"github.com/Arihantawasthi/sage.git/internal/logger"
	"github.com/Arihantawasthi/sage.git/internal/manager"
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

    processStore := manager.NewProcessStore(config)
    spmpServer := spmp.NewSPMPServer(config, logger, processStore)

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    go func() {
        if err := spmpServer.Start(); err != nil {
            fmt.Fprintf(os.Stderr, "SPMP server failed: %v", err)
            os.Exit(1)
        }
    }()

    <-ctx.Done()
}
