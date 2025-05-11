package spmp

import (
	"fmt"
	"net"

	"github.com/Arihantawasthi/sage.git/internal/models"
)

type connWriter struct {
	Conn net.Conn
}

type SPMPWriter interface {
	Write(en string, msgType byte, payload []byte) error
}

func NewSPMPWriter(conn net.Conn) *connWriter {
	return &connWriter{
		Conn: conn,
	}
}

func (w *connWriter) Write(en string, msgType byte, payload []byte) error {
	rspPkt, err := NewPacket(V1, en, msgType, payload)
	if err != nil {
		return err
	}

	data, err := rspPkt.Encode()
	if err != nil {
		return err
	}

	_, err = w.Conn.Write(data)
	return err
}

type SPMPRequest struct {
	Conn   net.Conn
	Packet *Packet
}

func NewSPMPRequest(conn net.Conn) (*SPMPRequest, error) {
	pkt, err := DecodePacket(conn)
	if err != nil {
		return nil, fmt.Errorf("decoding failed: %w", err)
	}

	return &SPMPRequest{
		Conn:   conn,
		Packet: pkt,
	}, nil
}

type HandlerFunc func(*SPMPRequest, SPMPWriter) error

type CommandMux struct {
	Cfg      models.Config
	handlers map[byte]HandlerFunc
}

func NewCommandMux() *CommandMux {
	return &CommandMux{
		handlers: make(map[byte]HandlerFunc),
	}
}

func (m *CommandMux) HandleCommand(cmdType byte, handler HandlerFunc) {
	m.handlers[cmdType] = handler
}

func (m *CommandMux) Serve(conn net.Conn) {
	defer conn.Close()

	req, err := NewSPMPRequest(conn)
	if err != nil {
		fmt.Printf("error creating request: %s", err)
	}
	handler, ok := m.handlers[req.Packet.Type]
	if !ok {
		fmt.Printf("no handler for command type: %v\n", req.Packet.Type)
		return
	}

	writer := &connWriter{Conn: conn}
	handler(req, writer)
}
