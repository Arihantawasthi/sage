package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Arihantawasthi/sage.git/internal/models"
)

func LoadConfig() (models.Config, error) {
    path := "../config.json"
	b, err := os.ReadFile(path)
	if err != nil {
        return models.Config{}, fmt.Errorf("error reading config file '%s': %w", path, err)
	}

	var config models.Config
    err = json.Unmarshal(b, &config)
    if err != nil {
        return models.Config{}, fmt.Errorf("error unmarshalling config file '%s': %w", path, err)
    }

	return config, nil
}
