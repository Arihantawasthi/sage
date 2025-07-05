package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Arihantawasthi/sage.git/internal/models"
)

func LoadConfig() (models.Config, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return models.Config{}, fmt.Errorf("error locating home directory")
    }

    confFilePath := fmt.Sprintf("%s/.sage/sage-conf.json", homeDir)
    _, err = os.Stat(confFilePath)
    if err != nil {
        os.OpenFile(confFilePath, os.O_CREATE, 0644)
        return models.Config{}, nil
    }

	b, err := os.ReadFile(confFilePath)
	if err != nil {
		return models.Config{}, fmt.Errorf("error reading config file '%s': %w", confFilePath, err)
	}

	var services models.Services
	err = json.Unmarshal(b, &services)
	if err != nil {
		return models.Config{}, fmt.Errorf("error unmarshalling config file '%s': %w", confFilePath, err)
	}

    m := make(map[string]models.Service)
    for _, svc := range services.Services {
        m[svc.Name] = svc
    }

	return models.Config{
        ServiceMap: m,
    }, nil
}
