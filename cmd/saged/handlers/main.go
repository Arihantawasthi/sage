package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/Arihantawasthi/sage.git/internal/spmp"
)

func GetListOfServices(r *spmp.SPMPRequest, w spmp.SPMPWriter) error {
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
