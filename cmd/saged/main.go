package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/shirou/gopsutil/mem"
)

type Services struct {
	Name    string   `json:"name"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type Config struct {
	Services []Services `json:"services"`
}

type Process struct {
	Name string
	Cmd  string
}

type Processes struct {
    Pid int
    Procs []Process
}

func loadConfig() (Config, error) {
	b, err := os.ReadFile("./config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading config file: %s", err)
	}
	var config Config
	json.Unmarshal(b, &config)

	return config, nil
}

func startServices(config Config) {
    var procs []Process
	for _, service := range config.Services {
		cmd := exec.Command(service.Command, service.Args...)
		fmt.Println(cmd)
		err := cmd.Start()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error starting services %s: %v", service.Name, err)
			continue
		}
        args := service.Args

        proc := Process{
            Name: service.Name,
            Cmd: fmt.Sprintf("%s %s", service.Command, args),
        }

        procs := Processes{
            Pid: cmd.Process.Pid,
            Procs: append(procs, proc),
        }
        fmt.Printf("%v\n", procs)
	}
    vm, err := mem.VirtualMemory()
    if err != nil {
        fmt.Fprintf(os.Stderr, "error getting memory info: %s", err)
    }
    fmt.Println(vm.Total)
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading config file: %s", err)
	}
	startServices(config)
}
