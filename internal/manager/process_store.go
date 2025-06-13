package manager

import (
	"fmt"
	"sync"

	"github.com/Arihantawasthi/sage.git/internal/models"
)

type ProcessStore struct {
	mu       sync.RWMutex
	cfg      models.Config
	services map[string]*models.Process
}

func NewProcessStore(config models.Config) *ProcessStore {
	return &ProcessStore{
		cfg:      config,
		services: make(map[string]*models.Process),
	}
}

func (ps *ProcessStore) StartProcess(serviceName string) {
    fmt.Println(ps.cfg.ServiceMap)
}

func (ps *ProcessStore) Get(serviceName string) {
}
