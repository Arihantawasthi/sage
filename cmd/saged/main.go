package main

import (
	"fmt"
	"os"

	"github.com/Arihantawasthi/sage.git/internal/config"
	"github.com/Arihantawasthi/sage.git/internal/logger"
	"github.com/Arihantawasthi/sage.git/internal/models"
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
    fmt.Println(logger, config)
}
