package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Arihantawasthi/sage.git/internal/models"
	"github.com/Arihantawasthi/sage.git/internal/spmp"
	"github.com/Arihantawasthi/sage.git/internal/utils"
)

func main() {
    command := os.Args[1]
    var cmdType byte
    if command == "list" {
        cmdType = spmp.TypeList
    }
    if command == "start" {
        cmdType = spmp.TypeStart
    }
    if command == "stop" {
        cmdType = spmp.TypeStop
    }

    serviceName := ""
    if len(os.Args) > 2 {
        serviceName = os.Args[2]
    }

    client := spmp.NewSPMPClient()
    packet, err := spmp.NewPacket(spmp.V1, spmp.TEXTEncoding, cmdType, []byte(serviceName))
    if err != nil {
        fmt.Fprintf(os.Stderr, "error while creating packet: %s", err)
    }

    receivedPkt, err := client.SendPacket(packet)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error while fetching info: %s", err)
    }

    var response models.Response[[]models.PListData]
    json.Unmarshal(receivedPkt.Payload, &response)
    if response.RequestStatus == 0 {
        fmt.Fprintf(os.Stderr, "%s\n", response.Msg)
        return
    }

    fmt.Fprintf(os.Stdout, "%s\n", response.Msg)
    if len(response.Data) > 0 {
        utils.PrintTable(response.Data)
    }
    return
}
