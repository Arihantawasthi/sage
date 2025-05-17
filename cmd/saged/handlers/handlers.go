package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/Arihantawasthi/sage.git/internal/logger"
	"github.com/Arihantawasthi/sage.git/internal/models"
	"github.com/Arihantawasthi/sage.git/internal/spmp"
)

type Handler struct {
	cfg    models.Config
	logger logger.SlogLogger
}

func NewHandler(config models.Config, logger logger.SlogLogger) *Handler {
    return &Handler{
        cfg: config,
        logger: logger,
    }
}

func (h *Handler) HandleListServices(r *spmp.SPMPRequest, w spmp.SPMPWriter) error {
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

func (h *Handler) HandleStartService(r *spmp.SPMPRequest, w spmp.SPMPWriter) error {
	if string(r.Packet.Encoding[:]) != spmp.TEXTEncoding {
        h.logger.Error("Invalid packet encoding", "func: HandleStartService", "START", "", "sagectl", string(r.Packet.Encoding[:]))
		resp := models.Response[uint8]{
			RequestStatus: 0,
			Msg:           "Execution failed, expected encoding is TEXT",
			Data:          0,
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
            h.logger.Error("error while marshalling response struct: ", "func: HandleStartService", "START", "", "sagectl", string(r.Packet.Encoding[:]))
            return nil
		}
		w.Write(spmp.JSONEncoding, spmp.TypeStart, respBytes)
        return nil
	}

	serviceName := string(r.Packet.Payload)

	serviceManager := NewProcessManager(h.cfg)
	resp, err := serviceManager.StartService(serviceName)
	if err != nil {
        h.logger.Error("error while starting service: ", "func: HandleStartService", "START", "", "sagectl", string(r.Packet.Encoding[:]))
	}

	respByte, err := json.Marshal(resp)
	if err != nil {
        h.logger.Error("error marshalling the repsonse: ", "func: HandleStartService", "START", "", "sagectl", string(r.Packet.Encoding[:]))
	}
	w.Write(spmp.JSONEncoding, spmp.TypeStart, respByte)
	return nil
}
