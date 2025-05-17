package spmp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	V1 byte = 0x01
	V2 byte = 0x02

	JSONEncoding string = "JS"
	TEXTEncoding string = "TX"

	TypeList   byte = 0x01
	TypeStatus byte = 0x02
	TypeStart  byte = 0x03
	TypeStop   byte = 0x04

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

type Payload struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func NewPacket(v byte, en string, msgType byte, payload []byte) (*Packet, error) {
	payloadLength := len(payload)
	if len(en) != 2 {
		return nil, fmt.Errorf("invalid encoding '%s'\n", en)
	}

	return &Packet{
		MagicBytes:  [2]byte{'S', 'G'},
		Version:     v,
		Encoding:    [2]byte{en[0], en[1]},
		Type:        msgType,
		PayloadSize: uint32(payloadLength),
		Payload:     payload,
	}, nil
}

func (p *Packet) Encode() ([]byte, error) {
	if int(p.PayloadSize) != len(p.Payload) {
		return nil, fmt.Errorf("payload size mismatch: declared: %d, actual: %d", p.PayloadSize, len(p.Payload))
	}

	buf := bytes.NewBuffer(make([]byte, 0, HeaderSize+p.PayloadSize))

	if err := binary.Write(buf, binary.BigEndian, p.MagicBytes); err != nil {
		return nil, fmt.Errorf("failed to write magic bytes: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, p.Version); err != nil {
		return nil, fmt.Errorf("failed to write version byte: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, p.Encoding); err != nil {
		return nil, fmt.Errorf("failed to write encoding bytes: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, p.Type); err != nil {
		return nil, fmt.Errorf("failed to write type bytes: %w", err)
	}
	if err := binary.Write(buf, binary.BigEndian, p.PayloadSize); err != nil {
		return nil, fmt.Errorf("failed to write payload size bytes: %w", err)
	}

	if _, err := buf.Write(p.Payload); err != nil {
		return nil, fmt.Errorf("failed to write payload bytes: %w", err)
	}

	return buf.Bytes(), nil
}

func DecodePacket(conn net.Conn) (*Packet, error) {
	headerBuf := make([]byte, HeaderSize)
	if _, err := io.ReadFull(conn, headerBuf); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	payloadSize := binary.BigEndian.Uint32(headerBuf[6:10])
	payloadBuf := make([]byte, payloadSize)
	if _, err := io.ReadFull(conn, payloadBuf); err != nil {
		return nil, fmt.Errorf("failed to read payload size: %w", err)
	}

	data := append(headerBuf, payloadBuf...)
	reader := bytes.NewReader(data)
	pkt := Packet{}

	if err := binary.Read(reader, binary.BigEndian, &pkt.MagicBytes); err != nil {
		return nil, fmt.Errorf("failed to read magic bytes: %w", err)
	}
	if pkt.MagicBytes != [2]byte{'S', 'G'} {
		return nil, fmt.Errorf("invalid magic bytes '%s'", pkt.MagicBytes[:])
	}

	if err := binary.Read(reader, binary.BigEndian, &pkt.Version); err != nil {
		return nil, fmt.Errorf("failed to read version byte: %w", err)
	}
	if pkt.Version != V1 {
		return nil, fmt.Errorf("invalid protocol verion: '%b'", pkt.Version)
	}

	if err := binary.Read(reader, binary.BigEndian, &pkt.Encoding); err != nil {
		return nil, fmt.Errorf("failed to read encoding bytes: %w", err)
	}
	if pkt.Encoding != [2]byte{JSONEncoding[0], JSONEncoding[1]} && pkt.Encoding != [2]byte{TEXTEncoding[0], TEXTEncoding[1]} {
		return nil, fmt.Errorf("invalid protocol verion: '%b'", pkt.Version)
	}

	if err := binary.Read(reader, binary.BigEndian, &pkt.Type); err != nil {
		return nil, fmt.Errorf("failed to read type bytes: %w", err)
	}
	if pkt.Type != TypeList && pkt.Type != TypeStatus && pkt.Type != TypeStart && pkt.Type != TypeStop {
		return nil, fmt.Errorf("invalid type: '%b'", pkt.Version)
	}

	if err := binary.Read(reader, binary.BigEndian, &pkt.PayloadSize); err != nil {
		return nil, fmt.Errorf("failed to read payload size bytes: %w", err)
	}

	if uint32(len(data)) < HeaderSize+pkt.PayloadSize {
		return nil, fmt.Errorf("payload size mismatch: declared %d, but only %d bytes were received", pkt.PayloadSize, uint32(len(data))-HeaderSize)
	}

    if pkt.PayloadSize > 0 {
        pkt.Payload = make([]byte, pkt.PayloadSize)
        if _, err := reader.Read(pkt.Payload); err != nil {
            return nil, fmt.Errorf("failed to read payload: %w", err)
        }
    } else {
        pkt.Payload = []byte{}
    }

	return &pkt, nil
}
