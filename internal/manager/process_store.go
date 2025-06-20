package manager

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/Arihantawasthi/sage.git/internal/models"
	"github.com/Arihantawasthi/sage.git/internal/utils"
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
	cmd := exec.Command(service.Command, service.Args...)
    cmd.Env = os.Environ()
    cmd.Dir = service.WorkingDir
    for k, v := range ps.cfg.ServiceMap[serviceName].Env {
        envVar := fmt.Sprintf("%s=%s", k, v)
        cmd.Env = append(cmd.Env, envVar)
    }

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return fmt.Sprintf("failed to get stdout pipe: %v", err)
    }
    stderr, err := cmd.StderrPipe()
    if err != nil {
        return fmt.Sprintf("failed to get stderr pipe: %v", err)
    }

    lDir, err := utils.CreateServiceLogDir()
    serviceLogPath := fmt.Sprintf("%s/%s.log", lDir, serviceName)
    if err != nil {
        e := fmt.Errorf("error starting the process: %v\n" ,err)
        return e.Error()
    }
    if _, err := os.Stat(serviceLogPath); err == nil {
        _, err := os.Create(serviceLogPath)
        if err != nil {
            e := fmt.Errorf("error creating log file for %s", serviceName)
            return e.Error()
        }
    }

    go utils.StreamLogs(stdout, fmt.Sprintf("[stdout][%s]", serviceName), serviceLogPath)
    go utils.StreamLogs(stderr, fmt.Sprintf("[stderr][%s]", serviceName), serviceLogPath)

	err = cmd.Start()
	if err != nil {
		e := fmt.Errorf("error starting the process: %v\n", err)
		return e.Error()
	}

	pid := cmd.Process.Pid

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

    time.Sleep(10 * time.Second)
    message := fmt.Sprintf("Service '%s' started successfully with PID %d\n", serviceName, pid)
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

		case <-stopChan:
			proc.Kill()
			return
		}
	}
}
