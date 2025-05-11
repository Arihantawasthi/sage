package handlers

import (
	"encoding/json"
	"fmt"
	"os"

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

	serviceManager := NewProcessManager(h.cfg)
	resp, err := serviceManager.StartService(serviceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while starting service: %s\n", err)
	}

	respByte, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling the reponse from StartService: %s\n", err)
	}
	w.Write(spmp.JSONEncoding, spmp.TypeStart, respByte)
	return nil
}
