package artnet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

// Decode NUL-terminated/padded string
func decodeString(buf []byte) string {
	for i := 0; i < len(buf); i++ {
		if buf[i] == 0 {
			return string(buf[:i])
		}
	}
	return string(buf)
}

type Transport struct {
	udpConn *net.UDPConn // listening on port
}

func (t *Transport) recv() (ArtPacket, *net.UDPAddr, error) {
	var header ArtHeader
	var buf = make([]byte, MTU)
	var srcAddr *net.UDPAddr

	if read, udpAddr, err := t.udpConn.ReadFromUDP(buf); err != nil {
		return nil, nil, fmt.Errorf("net:UDPConn.ReadFromUDP: %v", err)
	} else {
		buf = buf[:read]
		srcAddr = udpAddr
	}

	if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &header); err != nil {
		return nil, srcAddr, fmt.Errorf("binary.Read ArtHeader: %v", err)
	}

	if header.ID != ARTNET {
		return nil, srcAddr, fmt.Errorf("Invalid magic")
	}

	if packet, err := t.decode(header, buf); err != nil {
		return nil, srcAddr, err
	} else {
		return packet, srcAddr, err
	}
}

func (t *Transport) decode(header ArtHeader, buf []byte) (ArtPacket, error) {
	var packet ArtPacket

	switch header.OpCode {
	case OpPoll:
		packet = &ArtPoll{
			ArtHeader: header,
		}

	case OpPollReply:
		packet = &ArtPollReply{
			ArtHeader: header,
		}

	case OpDmx:
		packet = &ArtDmx{
			ArtHeader: header,
		}

	default:
		return nil, fmt.Errorf("Unknown opcode: %04x", header.OpCode)
	}

	if err := binary.Read(bytes.NewReader(buf), binary.BigEndian, packet); err != nil {
		return nil, fmt.Errorf("binary.Read %T: %v", packet, err)
	}

	return packet, nil
}

func (t *Transport) send(addr *net.UDPAddr, packet ArtPacket, data []byte) error {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, packet); err != nil {
		return err
	}

	if _, err := buf.Write(data); err != nil {
		return err
	}

	// t.log.Debugf("send: opcode=%04x len=%v", packet.Header().OpCode, buf.Len())

	if _, err := t.udpConn.WriteToUDP(buf.Bytes(), addr); err != nil {
		return err
	}

	return nil
}
