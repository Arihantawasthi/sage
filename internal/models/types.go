package models

type Service struct {
	Name    string   `json:"name"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
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

type Process struct {
	Pid        int32
	PName      string
	Name       string
	Cmd        string
	UpTime     string
	CPUPercent float64
	MemPrecent float32
	StopChan   chan struct{}
}
