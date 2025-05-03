package spmp

import (
	"fmt"
	"net"
)

type SPMPClient struct {
    socketPath string
}

func NewSPMPClient() *SPMPClient {
    return &SPMPClient{
        socketPath: "/tmp/sage.sock",
    }
}

func (c *SPMPClient) Show() (string, error) {
    conn, err := net.Dial("unix", c.socketPath)
    if err != nil {
        return "", fmt.Errorf("error connecting to unix socket: %w", err)
    }
    defer conn.Close()
    message := []byte("hello")
    _, err = conn.Write(message)
    if err != nil {
        return "", fmt.Errorf("error writing to the connection: %v\n", err)
    }

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("error reading from the connection: %v\n", err)
	}
    return string(buffer[:n]), nil
}
