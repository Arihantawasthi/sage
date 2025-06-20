package spmp

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/Arihantawasthi/sage.git/internal/logger"
	"github.com/Arihantawasthi/sage.git/internal/manager"
	"github.com/Arihantawasthi/sage.git/internal/models"
)

type SPMPServer struct {
	cfg    models.Config
	logger *logger.SlogLogger
	router map[byte]func(*Packet) ([]byte, string, error)
	ps     *manager.ProcessStore
}

func NewSPMPServer(cfg models.Config, logger *logger.SlogLogger, processStore *manager.ProcessStore) *SPMPServer {
	s := &SPMPServer{
		cfg:    cfg,
		logger: logger,
		router: make(map[byte]func(*Packet) ([]byte, string, error)),
		ps:     processStore,
	}
	s.router[TypeStart] = s.handleStart
	s.router[TypeList] = s.handleList
	s.router[TypeStop] = s.handleStop

	return s
}

func (s *SPMPServer) Start() error {
    socketPath := "/tmp/sage.sock"
    if _, err := os.Stat(socketPath); err == nil {
        err := os.Remove(socketPath)
        if err != nil {
            s.logger.Error("Failed to remove existing socket", "START", "os.Remove", "", "", "")
            return err
        }
        s.logger.Info("Removed existing socket file", "START", "os.Remove", "", "", "")
    }

	listener, err := net.Listen("unix", socketPath)
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

	responsePayload, encoding, err := handler(packet)
	if err != nil {
		s.logger.Error("Handler failed", "func: handleConnection", "handler", remoteAddr, "", err)
		return
	}

	responsePkt, err := NewPacket(V1, encoding, packet.Type, responsePayload)
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

func (s *SPMPServer) handleStart(pkt *Packet) ([]byte, string, error) {
	serviceName := string(pkt.Payload)
	_, exists := s.cfg.ServiceMap[serviceName]
	if !exists {
		e := fmt.Sprintf("'%s': service name doesn't exist", serviceName)
		return []byte(e), TEXTEncoding, nil
	}
	message := s.ps.StartProcess(serviceName)
	return []byte(message), TEXTEncoding, nil
}

func (s *SPMPServer) handleStop(pkt *Packet) ([]byte, string, error) {
	serviceName := string(pkt.Payload)
	_, exists := s.cfg.ServiceMap[serviceName]
	if !exists {
		e := fmt.Sprintf("'%s': service name doesn't exist", serviceName)
		return []byte(e), TEXTEncoding, nil
	}
	message := s.ps.StopProcess(serviceName)
	return []byte(message), TEXTEncoding, nil
}

func (s *SPMPServer) handleList(pkt *Packet) ([]byte, string, error) {
	payload := string(pkt.Payload)
	plistData := s.ps.ListProcesses(payload)
	data := models.Response[[]models.PListData]{
		RequestStatus: 1,
		Msg:           "Service list retrieved successfully",
		Data:          plistData,
	}
	response, err := json.Marshal(data)
	if err != nil {
		e := fmt.Errorf("error in encoding json: %v", err)
		return []byte(e.Error()), TEXTEncoding, nil
	}

	return response, JSONEncoding, nil
}
