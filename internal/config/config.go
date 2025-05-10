package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Arihantawasthi/sage.git/internal/models"
)

func LoadConfig() (models.Config, error) {
	path := "./config.json"
	b, err := os.ReadFile(path)
	if err != nil {
		return models.Config{}, fmt.Errorf("error reading config file '%s': %w", path, err)
	}

	var services models.Services
	err = json.Unmarshal(b, &services)
	if err != nil {
		return models.Config{}, fmt.Errorf("error unmarshalling config file '%s': %w", path, err)
	}

    m := make(map[string]models.Service)
    for _, svc := range services.Services {
        m[svc.Name] = svc
    }

	return models.Config{
        ServiceMap: m,
    }, nil
}
