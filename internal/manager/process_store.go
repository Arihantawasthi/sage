package manager

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/Arihantawasthi/sage.git/internal/models"
	"github.com/shirou/gopsutil/process"
)

type ProcessStore struct {
	mu    sync.RWMutex
	cfg   models.Config
	store map[string]*models.Process
}

func NewProcessStore(config models.Config) *ProcessStore {
	return &ProcessStore{
		cfg:   config,
		store: make(map[string]*models.Process),
	}
}

func (ps *ProcessStore) StartProcess(serviceName string) string {
	service := ps.cfg.ServiceMap[serviceName]
	fmt.Println(service.Command, service.Args)
	cmd := exec.Command(service.Command, service.Args...)
	err := cmd.Start()
	if err != nil {
		e := fmt.Errorf("error starting the process: %v\n", err)
		return e.Error()
	}

	pid := cmd.Process.Pid
	fmt.Println("Process ID: ", pid)

	ps.mu.Lock()
	stopChan := make(chan struct{})
	ps.store[serviceName] = &models.Process{
		Pid:      pid,
		Name:     serviceName,
		PName:    serviceName,
		Cmd:      service.Command,
		StopChan: stopChan,
	}
	ps.mu.Unlock()

	go ps.monitorProcess(serviceName, pid, stopChan)

	go func() {
		err := cmd.Wait()
		delete(ps.store, serviceName)
		if err != nil {
			fmt.Printf("process %s (PID %d) exited with error: %v\n", serviceName, pid, err)
		} else {
			fmt.Printf("process %s (PID %d) exited normally\n", serviceName, pid)
		}
	}()

	message := fmt.Sprintf("Service '%s' started successfully with PID %d", serviceName, pid)
	return message
}

func (ps *ProcessStore) StopProcess(serviceName string) string {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	runningProcess, exists := ps.store[serviceName]
	if !exists {
		e := fmt.Errorf("Service %s is not running", serviceName)
		return e.Error()
	}
	close(runningProcess.StopChan)
	delete(ps.store, serviceName)

	message := fmt.Sprintf("Service '%s' stopped successfully", serviceName)
	return message
}

func (ps *ProcessStore) ListProcesses(payload string) []models.PListData {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var plist []models.PListData
    for name, service := range ps.cfg.ServiceMap {
		rp, exists := ps.store[name]
		data := models.PListData{
			Pid:        0,
			PName:      service.Name,
			Name:       service.Name,
			Cmd:        service.Command,
			Status:     "offline",
			UpTime:     "0s",
			CPUPercent: 0.00,
			MemPrecent: 0.00,
		}

		if exists && rp != nil {
			data.Pid = rp.Pid
			data.Status = "online"
			data.UpTime = rp.UpTime
			data.CPUPercent = rp.CPUPercent
			data.MemPrecent = rp.MemPrecent
		}

		plist = append(plist, data)
	}

	return plist
}

func (ps *ProcessStore) monitorProcess(serviceName string, pid int, stopChan chan struct{}) {
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		fmt.Errorf("failed to create process monitor for PID %d: %v", pid, err)
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cpuPercent, err := proc.CPUPercent()
			if err != nil {
				fmt.Errorf("process %s (PID %d) is not running\n", serviceName, pid)
				continue
			}
			memPercent, err := proc.MemoryPercent()
			if err != nil {
				fmt.Errorf("process %s (PID %d) is not running\n", serviceName, pid)
				continue
			}
			createTimeMillis, err := proc.CreateTime()
			if err != nil {
				fmt.Errorf("failed to get the create time for process %s: %v\n", serviceName, pid)
			}
			startTime := time.Unix(0, createTimeMillis*int64(time.Millisecond))
			uptime := time.Since(startTime).String()

			ps.mu.Lock()
			storedProc, exists := ps.store[serviceName]
			if exists {
				storedProc.CPUPercent = cpuPercent
				storedProc.MemPrecent = memPercent
				storedProc.UpTime = uptime
			}
			ps.mu.Unlock()
			fmt.Println("Process: ", proc.Pid, cpuPercent, memPercent, uptime)

		case <-stopChan:
			proc.Kill()
			return
		}
	}
}
