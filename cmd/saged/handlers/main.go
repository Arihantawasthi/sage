package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Arihantawasthi/sage.git/internal/models"
	"github.com/Arihantawasthi/sage.git/internal/spmp"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

type ProcessManager struct {
    Processes map[string]models.Process
    Mutex sync.Mutex
    Cfg models.Config
}

func NewProcessManager(cfg models.Config) *ProcessManager{
    return &ProcessManager{
        Processes: make(map[string]models.Process),
        Cfg: cfg,
    }
}

func GetListOfServices(r *spmp.SPMPRequest, w spmp.SPMPWriter, cfg models.Config) error {
	payload := spmp.Payload{
		Name: "gitbook",
		Type: "list",
	}
	fmt.Println(payload)
	payloadBytes, err := json.Marshal(payload)
	fmt.Println(payloadBytes)
	if err != nil {
		return fmt.Errorf("error encoding json into bytes: %v", err)
	}
	w.Write(spmp.JSONEncoding, spmp.TypeStatus, payloadBytes)
	return nil
}

func HandleStartService(r *spmp.SPMPRequest, w spmp.SPMPWriter, cfg models.Config) error {
	if string(r.Packet.Encoding[:]) == spmp.TEXTEncoding {
		resp := models.Response[uint8]{
			RequestStatus: 0,
			Msg:           "Execution failed, expected encoding is TEXT",
			Data:          0,
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s", err)
		}
		fmt.Fprintf(os.Stderr, "Error: Expected encoding, TEXT")
		w.Write(spmp.JSONEncoding, spmp.TypeStart, respBytes)
	}

    serviceName := string(r.Packet.Payload)

    serviceManager := NewProcessManager(cfg)
    if serviceName == "all" {
        serviceManager.StartServices()
    }

    go func() {
    }()
	return nil
}

func (p *ProcessManager) StartServices() {
	var procs []models.Process
	for _, service := range p.Cfg.Services {
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
}
