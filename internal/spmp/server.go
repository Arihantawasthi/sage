package spmp

import (
	"fmt"
	"net"

	"github.com/Arihantawasthi/sage.git/internal/logger"
	"github.com/Arihantawasthi/sage.git/internal/manager"
	"github.com/Arihantawasthi/sage.git/internal/models"
)

type SPMPServer struct {
	cfg    models.Config
	logger *logger.SlogLogger
	router map[byte]func(*Packet) ([]byte, error)
	ps     *manager.ProcessStore
}

func NewSPMPServer(cfg models.Config, logger *logger.SlogLogger, processStore *manager.ProcessStore) *SPMPServer {
	s := &SPMPServer{
		cfg:    cfg,
		logger: logger,
		router: make(map[byte]func(*Packet) ([]byte, error)),
		ps:     processStore,
	}
	s.router[TypeStart] = s.handleStart
	s.router[TypeList] = s.handleStart
    s.router[TypeStop] = s.handleStop

	return s
}

func (s *SPMPServer) Start() error {
	listener, err := net.Listen("unix", "/tmp/sage.sock")
	if err != nil {
		s.logger.Error("Failed to start SPMP server", "START", "net.Listen", "", "", "")
		return err
	}
	defer listener.Close()

	s.logger.Info("SPMP Server Started", "START", "net.Listen", "", "", "")

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error("Failed to accept connection", "START", "Accept", "", "", "")
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *SPMPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	s.logger.Info("Accepted new connection", "func: handleConnection", "Accept", remoteAddr, "", "")

	packet, err := DecodePacket(conn)
	if err != nil {
		s.logger.Info("Error in decoding the packet", "func: handleConnection", "DecodePacket", remoteAddr, "", err)
		return
	}

	handler, exists := s.router[packet.Type]
	if !exists {
		err = fmt.Errorf("unknown command type: %d", packet.Type)
		s.logger.Error("No handler found for command type", "func: handleConnection", "router lookup", remoteAddr, "", err)
		return
	}

	responsePayload, err := handler(packet)
	if err != nil {
		s.logger.Error("Handler failed", "func: handleConnection", "handler", remoteAddr, "", err)
		return
	}

	responsePkt, err := NewPacket(V1, TEXTEncoding, packet.Type, responsePayload)
	if err != nil {
		s.logger.Error("failed to build response packet", "func: handleConnection", "NewPacket", remoteAddr, "", err)
		return
	}

	encoded, err := responsePkt.Encode()
	if err != nil {
		s.logger.Error("failed to encode packet", "func: handleConnection", "Encode", remoteAddr, "", err)
		return
	}

	if _, err := conn.Write(encoded); err != nil {
		s.logger.Error("Failed to write the response", "func: handleConnection", "Write", remoteAddr, "", err)
		return
	}
}

func (s *SPMPServer) handleStart(pkt *Packet) ([]byte, error) {
	serviceName := string(pkt.Payload)
	_, exists := s.cfg.ServiceMap[serviceName]
	if !exists {
		e := fmt.Sprintf("'%s': service name doesn't exist", serviceName)
		return []byte(e), nil
	}
	message := s.ps.StartProcess(serviceName)
	return []byte(message), nil
}

func (s *SPMPServer) handleStop(pkt *Packet) ([]byte, error) {
	serviceName := string(pkt.Payload)
	_, exists := s.cfg.ServiceMap[serviceName]
	if !exists {
		e := fmt.Sprintf("'%s': service name doesn't exist", serviceName)
		return []byte(e), nil
	}
	message := s.ps.StopProcess(serviceName)
	return []byte(message), nil
}
