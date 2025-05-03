package main

import (
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

    packetBytes, err := packet.Encode()
    if err != nil {
        fmt.Fprintf(os.Stderr, "error while encoding packet: %s", err)
    }

    data, err := client.SendPacket(packetBytes)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error while fetching info: %s", err)
    }

    fmt.Println(data)
}
