package spmp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
    V1 byte = 0x01
    V2 byte = 0x02

    JSONEncoding string = "JS"
    TEXTEncoding string = "TX"

    TypeList byte = 0x01
    TypeStatus byte = 0x02
    TypeStart byte = 0x03
    TypeStop byte = 0x04

    HeaderSize uint32 = 10
)

type Packet struct {
	MagicBytes  [2]byte
	Version     byte
	Encoding    [2]byte
	Type        byte
	PayloadSize uint32
	Payload     []byte
}

func NewPacket(v byte, en string, msgType byte, payload []byte) (*Packet, error) {
    payloadLength := len(payload)
    if len(en) != 2 {
        return nil, fmt.Errorf("invalid encoding '%s'\n", en)
    }

    return &Packet{
        MagicBytes: [2]byte{'S', 'G'},
        Version: v,
        Encoding: [2]byte{en[0], en[1]},
        Type: msgType,
        PayloadSize: uint32(payloadLength),
        Payload: payload,
    }, nil
}

func (p *Packet) Encode() ([]byte, error) {
    if int(p.PayloadSize) != len(p.Payload) {
        return nil, fmt.Errorf("payload size mismatch: declared: %d, actual: %d", p.PayloadSize, len(p.Payload))
    }

    buf := bytes.NewBuffer(make([]byte, 0, HeaderSize + p.PayloadSize))

    if err := binary.Write(buf, binary.BigEndian, p.MagicBytes); err != nil {
        return nil, err
    }
    if err := binary.Write(buf, binary.BigEndian, p.Version); err != nil {
        return nil, err
    }
    if err := binary.Write(buf, binary.BigEndian, p.Encoding); err != nil {
        return nil, err
    }
    if err := binary.Write(buf, binary.BigEndian, p.PayloadSize); err != nil {
        return nil, err
    }

    if _, err := buf.Write(p.Payload); err != nil {
        return nil, fmt.Errorf("writing payload: %w", err)
    }

    return buf.Bytes(), nil
}
