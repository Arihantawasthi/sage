package spmp

import (
	"fmt"
	"net"
)

type SPMPClient struct {
	socketPath string
}

func NewSPMPClient() *SPMPClient {
	return &SPMPClient{
		socketPath: "/tmp/sage.sock",
	}
}

func (c *SPMPClient) SendAndReceivePacket(pkt *Packet) (*Packet, error) {
	conn, err := net.Dial("unix", c.socketPath)
	if err != nil {
		return nil, fmt.Errorf("error connecting to unix socket: %w", err)
	}
	defer conn.Close()
    packetBytes, err := pkt.Encode()
    if err != nil {
        return nil, fmt.Errorf("error while ecoding packet: %s", err)
    }

	_, err = conn.Write(packetBytes)
	if err != nil {
		return nil, fmt.Errorf("error writing to the connection: %v\n", err)
	}

    decodedPkt, err := DecodePacket(conn)
    if err != nil {
        return nil, fmt.Errorf("error while decoding the received bytes: %s", err)
    }
    return decodedPkt, nil
}
