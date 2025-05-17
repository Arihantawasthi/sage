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
	sm     *ProcessManager
}

func NewHandler(config models.Config, logger logger.SlogLogger) *Handler {
	serviceManager := NewProcessManager(config)
	return &Handler{
		cfg:    config,
		logger: logger,
		sm:     serviceManager,
	}
}

func (h *Handler) HandleListServices(r *spmp.SPMPRequest, w spmp.SPMPWriter) error {
    resp, err := h.sm.ListServices()
    respBytes, err := json.Marshal(resp)
    if err != nil {
        errMsg := fmt.Sprintf("error while marshalling response struct: %s", err)
        h.logger.Error(errMsg, "func: HandleListServices", "LIST", "", "sagectl", string(r.Packet.Encoding[:]))
        return nil
    }

    w.Write(spmp.JSONEncoding, spmp.TypeStart, respBytes)
	return nil
}

func (h *Handler) HandleStopService(r *spmp.SPMPRequest, w spmp.SPMPWriter) error {
	if string(r.Packet.Encoding[:]) != spmp.TEXTEncoding {
		h.logger.Error("Invalid packet encoding", "func: HandleStopService", "STOP", "", "sagectl", string(r.Packet.Encoding[:]))
		resp := models.Response[string]{
			RequestStatus: 0,
			Msg:           "Execution failed, expected encoding is TEXT",
			Data:          "",
		}
		respBytes, err := json.Marshal(resp)
		if err != nil {
            errMsg := fmt.Sprintf("error while marshalling response struct: %s", err)
			h.logger.Error(errMsg, "func: HandleStopService", "STOP", "", "sagectl", string(r.Packet.Encoding[:]))
			return nil
		}
		w.Write(spmp.JSONEncoding, spmp.TypeStart, respBytes)
		return nil
	}

    serviceName := string(r.Packet.Payload)
    resp, err := h.sm.StopService(serviceName)
    if err != nil {
        errMsg := fmt.Sprintf("error while stopping the service: %s", err)
        h.logger.Error(errMsg, "func: HandleStopService", "STOP", "", "sagectl", string(r.Packet.Encoding[:]))
    }

    respByte, err := json.Marshal(resp)
    if err != nil {
        errMsg := fmt.Sprintf("error while marshalling response struct: %s", err)
		h.logger.Error(errMsg, "func: HandleStopService", "STOP", "", "sagectl", string(r.Packet.Encoding[:]))
    }

    w.Write(spmp.JSONEncoding, spmp.TypeStop, respByte)
    return nil
}

func (h *Handler) HandleStartService(r *spmp.SPMPRequest, w spmp.SPMPWriter) error {
	if string(r.Packet.Encoding[:]) != spmp.TEXTEncoding {
		h.logger.Error("Invalid packet encoding", "func: HandleStartService", "START", "", "sagectl", string(r.Packet.Encoding[:]))
		resp := models.Response[string]{
			RequestStatus: 0,
			Msg:           "Execution failed, expected encoding is TEXT",
			Data:          "",
        }
		respBytes, err := json.Marshal(resp)
		if err != nil {
            errMsg := fmt.Sprintf("error while marshalling response struct: %s", err)
			h.logger.Error(errMsg, "func: HandleStartService", "START", "", "sagectl", string(r.Packet.Encoding[:]))
			return nil
		}
		w.Write(spmp.JSONEncoding, spmp.TypeStart, respBytes)
		return nil
	}

	serviceName := string(r.Packet.Payload)

	resp, err := h.sm.StartService(serviceName)
	if err != nil {
        errMsg := fmt.Sprintf("error while starting service: %s", err)
		h.logger.Error(errMsg, "func: HandleStartService", "START", "", "sagectl", string(r.Packet.Encoding[:]))
	}

	respByte, err := json.Marshal(resp)
	if err != nil {
        errMsg := fmt.Sprintf("error while marshalling response struct: %s", err)
		h.logger.Error(errMsg, "func: HandleStartService", "START", "", "sagectl", string(r.Packet.Encoding[:]))
	}
	w.Write(spmp.JSONEncoding, spmp.TypeStart, respByte)
	return nil
}
