package models

type Services struct {
	Name    string `json:"name"`
	Command string `json:"command"`
	Args    []string `json:"args"`
}

type Config struct {
    Services []Services `json:"success"`
}
