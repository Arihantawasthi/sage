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
    if len(os.Args) < 2 {
        fmt.Fprintf(os.Stderr, "Usage: sagectl [list|status|start|stop] <service-name>\n")
        return
    }

    command := os.Args[1]
    serviceName := ""
    if len(os.Args) < 3 && command != "list" {
        fmt.Fprintf(os.Stderr, "Service name requried.\nExample Usage: sagectl %s <service-name>\n", command)
        return
    }
    if command != "list" {
        serviceName = os.Args[2]
    }

    client := spmp.NewSPMPClient()
    packet, err := buildPacket(command, serviceName)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v", err)
        os.Exit(1)
    }

    receivedPkt, err := client.SendAndReceivePacket(packet)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error while fetching info: %s\n", err)
        os.Exit(1)
    }

    if string(receivedPkt.Encoding[:]) == spmp.TEXTEncoding {
        fmt.Fprintf(os.Stdout, "%s\n", receivedPkt.Payload)
        return
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

func buildPacket(cmd, serviceName string) (*spmp.Packet, error) {
    var msgType byte
    switch cmd{
    case "list":
        msgType = spmp.TypeList
    case "start":
        msgType = spmp.TypeStart
    case "stop":
        msgType = spmp.TypeStop
    case "status":
        msgType = spmp.TypeStatus
    default:
        return nil, fmt.Errorf("unkown command: %s", cmd)
    }

    packet, err := spmp.NewPacket(spmp.V1, spmp.TEXTEncoding, msgType, []byte(serviceName))
    if err != nil {
        fmt.Fprintf(os.Stderr, "error while creating packet: %s\n", err)
    }

    return packet, err
}
