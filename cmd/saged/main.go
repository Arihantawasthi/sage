package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"

	"github.com/Arihantawasthi/sage.git/cmd/saged/handlers"
	"github.com/Arihantawasthi/sage.git/internal/config"
	"github.com/Arihantawasthi/sage.git/internal/models"
	"github.com/Arihantawasthi/sage.git/internal/spmp"
)

type Process struct {
	Pid        int32
	PName      string
	Name       string
	Cmd        string
	upTime     string
	CPUPercent float64
	MemPrecent float32
}

func startServices(config models.Config) {
	var procs []Process
	for _, service := range config.Services {
		cmd := exec.Command(service.Command, service.Args...)
		err := cmd.Start()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error starting services %s: %v", service.Name, err)
			continue
		}
		args := service.Args
		p, err := process.NewProcess(int32(cmd.Process.Pid))

		procCreateTime, err := p.CreateTime()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while getting create time for [%d]: %s", cmd.Process.Pid, err)
		}
		createTime := time.UnixMilli(procCreateTime)
		upTime := time.Since(createTime)

		procCPUPercent, err := p.CPUPercent()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while getting cpu usage for [%d]: %s", cmd.Process.Pid, err)
		}

		procMemPercent, err := p.MemoryPercent()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while getting memory usage for [%d]: %s", cmd.Process.Pid, err)
		}

		procName, err := p.Name()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while getting name for [%d]: %s", cmd.Process.Pid, err)
		}

		proc := Process{
			Pid:        p.Pid,
			PName:      procName,
			Name:       service.Name,
			Cmd:        fmt.Sprintf("%s %s", service.Command, strings.Join(args, "")),
			upTime:     upTime.String(),
			CPUPercent: procCPUPercent,
			MemPrecent: procMemPercent,
		}

		fmt.Println(proc)
		procs = append(procs, proc)
	}
	vm, err := mem.VirtualMemory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting memory info: %s", err)
	}
	fmt.Println(float64(vm.Total) / 1024 / 1024)
	fmt.Println(float64(vm.Used) / 1024 / 1024)
	fmt.Println(vm.UsedPercent)
}

func main() {
    config, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading config file: %s", err)
        os.Exit(1)
	}

    ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
    defer cancel()

    cmdMux := spmp.NewCommandMux()
    cmdMux.HandleCommand(spmp.TypeList, handlers.GetListOfServices)

    spmpServer := spmp.NewSPMPServer(config, cmdMux)
    go func(ctx context.Context) {
        spmpServer.ListenAndServe(ctx)
    }(ctx)

    <-ctx.Done()
    fmt.Println("Exiting...")
	// startServices(config)
}
