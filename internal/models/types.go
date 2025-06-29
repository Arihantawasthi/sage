package models

type Service struct {
	Name       string            `json:"name"`
	Command    string            `json:"command"`
	Args       []string          `json:"args"`
	WorkingDir string            `json:"workingDir"`
	Env        map[string]string `json:"env,omitempty"`
}

type Services struct {
	Services []Service `json:"services"`
}

type Config struct {
	ServiceMap map[string]Service `json:"serviceMap"`
}

type Response[T any] struct {
	RequestStatus uint8  `json:"requestStatus"`
	Msg           string `json:"msg"`
	Data          T      `json:"data"`
}

type PListData struct {
	Pid        int     `json:"pid"`
	PName      string  `json:"pname"`
	Name       string  `json:"name"`
	Cmd        string  `json:"cmd"`
	Status     string  `json:"status"`
	UpTime     string  `json:"uptime"`
	CPUPercent float64 `json:"cpuPercent"`
	MemPrecent float32 `json:"memPercent"`
}

type Process struct {
	Pid        int
	PName      string
	Name       string
	Cmd        string
	UpTime     string
	CPUPercent float64
	MemPrecent float32
	StopChan   chan struct{}
}

const LogFilePath = ".sage/saged/saged.log"
