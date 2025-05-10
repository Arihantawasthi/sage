package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Arihantawasthi/sage.git/internal/models"
	"github.com/Arihantawasthi/sage.git/internal/spmp"
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

    serviceName := ""
    if len(os.Args) > 1 {
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

    var response models.Response[interface{}]
    json.Unmarshal(receivedPkt.Payload, &response)
    fmt.Println(response)
}
