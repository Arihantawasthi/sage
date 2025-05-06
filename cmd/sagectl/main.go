package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Arihantawasthi/sage.git/internal/spmp"
)

func main() {
    command := os.Args[1]
    var cmdType byte
    if command == "list" {
        cmdType = spmp.TypeList
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

    var payload spmp.Payload
    json.Unmarshal(receivedPkt.Payload, &payload)
    fmt.Println(payload)
}
