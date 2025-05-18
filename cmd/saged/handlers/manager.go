package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Arihantawasthi/sage.git/internal/models"
	mem "github.com/shirou/gopsutil/mem"
	ps "github.com/shirou/gopsutil/process"
)

type ProcessManager struct {
	ProcessMap map[string]*models.Process
	Mutex      sync.Mutex
	Cfg        models.Config
}

func NewProcessManager(cfg models.Config) *ProcessManager {
	return &ProcessManager{
		ProcessMap: make(map[string]*models.Process),
		Cfg:        cfg,
	}
}

func (p *ProcessManager) ListServices() (models.Response[[]models.PListData], error) {
	var res []models.PListData
	for k, _ := range p.Cfg.ServiceMap {
		cfgPs, exists := p.ProcessMap[k]
		if !exists {
			res = append(res, models.PListData{
				Pid:        0,
				PName:      k,
				Name:       k,
				Cmd:        "-",
				Status:     "offline",
				UpTime:     "0s",
				CPUPercent: 0.0,
				MemPrecent: 0.0,
			})
		} else {
			res = append(res, models.PListData{
				Pid:        cfgPs.Pid,
				PName:      cfgPs.PName,
				Name:       cfgPs.Name,
				Cmd:        cfgPs.Cmd,
				Status:     "online",
				UpTime:     cfgPs.UpTime,
				CPUPercent: cfgPs.CPUPercent,
				MemPrecent: cfgPs.MemPrecent,
			})
		}
	}
	if len(res) == 0 {
		return models.Response[[]models.PListData]{
			RequestStatus: 1,
			Msg:           "No services running",
			Data:          res,
		}, nil
	}

	return models.Response[[]models.PListData]{
		RequestStatus: 1,
		Msg:           "List fetched successfully",
		Data:          res,
	}, nil
}

func (p *ProcessManager) StopService(name string) (models.Response[string], error) {
	_, exists := p.Cfg.ServiceMap[name]
	if !exists {
		return models.Response[string]{
			RequestStatus: 0,
			Msg:           "No services found with this name in the configuration",
			Data:          "",
		}, nil
	}

	runningPs, exists := p.ProcessMap[name]
	if !exists {
		return models.Response[string]{
			RequestStatus: 0,
			Msg:           "This service is not running",
			Data:          "",
		}, nil
	}

	close(runningPs.StopChan)
	delete(p.ProcessMap, name)

	response := models.Response[string]{
		RequestStatus: 1,
		Msg:           "Service stopped successfully",
		Data:          "",
	}
	return response, nil
}

func (p *ProcessManager) StartService(name string) (models.Response[string], error) {
	_, exists := p.ProcessMap[name]
	if exists {
		return models.Response[string]{
			RequestStatus: 0,
			Msg:           "Service already running",
			Data:          "",
		}, nil
	}

	service, exists := p.Cfg.ServiceMap[name]
	if !exists {
		return models.Response[string]{}, fmt.Errorf("No services found with this name in the configuration")
	}

	cmd := exec.Command(service.Command, service.Args...)
	err := cmd.Start()
	if err != nil {
		return models.Response[string]{}, fmt.Errorf("failed to start service: %v", err)
	}

	pid := int32(cmd.Process.Pid)
	proc := &models.Process{
		Pid:      pid,
		PName:    name,
		Name:     name,
		Cmd:      service.Command,
		StopChan: make(chan struct{}),
	}

	p.Mutex.Lock()
	p.ProcessMap[service.Name] = proc
	p.Mutex.Unlock()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		procHandle, err := ps.NewProcess(pid)
		if err != nil {
			fmt.Printf("Could not get process handle for PID %d: %v", pid, err)
			return
		}

		waitCh := make(chan error, 1)
		go func() {
			waitCh <- cmd.Wait()
		}()

		for {
			select {
			case <-proc.StopChan:
				fmt.Printf("Stop signal received for the service: %s", name)
				_ = procHandle.Kill()
				return

			case err := <-waitCh:
				if err != nil {
					if exitErr, ok := err.(*exec.ExitError); ok {
						fmt.Printf("Process [%s] exited with status: %d\n", name, exitErr.ExitCode())
					} else {
						fmt.Printf("Process [%s] exited with error: %s\n", name, err)
					}
				} else {
					fmt.Printf("Process [%s] exited normally\n", name)
				}

				p.Mutex.Lock()
				delete(p.ProcessMap, name)
				p.Mutex.Unlock()

				fmt.Printf("Cleaned up process state: %s\n", name)
				return

			case <-ticker.C:
				cpuPercent, cpuErr := procHandle.CPUPercent()
				memPercent, memErr := procHandle.MemoryPercent()
				createTime, timeErr := procHandle.CreateTime()

				if cpuErr == nil {
					proc.CPUPercent = cpuPercent
				}
				if memErr == nil {
					proc.MemPrecent = memPercent
				}
				if timeErr == nil {
					up := time.Since(time.UnixMilli(createTime)).Truncate(time.Second)
					proc.UpTime = up.String()
				}

				fmt.Printf("[%s] CPU: %.2f%% MEM: %.2f%% UPTIME: %s\n", name, proc.CPUPercent, proc.MemPrecent, proc.UpTime)
			}

		}
	}()

	return models.Response[string]{
		RequestStatus: 1,
		Msg:           "Service started!",
		Data:          fmt.Sprintf("PID: %d", pid),
	}, nil
}

func (p *ProcessManager) StartServices() error {
	var procs []models.Process
	for _, service := range p.Cfg.ServiceMap {
		cmd := exec.Command(service.Command, service.Args...)
		err := cmd.Start()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error starting services %s: %v", service.Name, err)
			continue
		}
		args := service.Args
		p, err := ps.NewProcess(int32(cmd.Process.Pid))

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

		proc := models.Process{
			Pid:        p.Pid,
			PName:      procName,
			Name:       service.Name,
			Cmd:        fmt.Sprintf("%s %s", service.Command, strings.Join(args, "")),
			UpTime:     upTime.String(),
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
	return nil
}
