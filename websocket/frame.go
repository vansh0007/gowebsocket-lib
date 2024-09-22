package websocket

import (
	"io"
)

const (
	// FinBit and Opcode types
	FinBit = 1 << 7

	// Opcodes for different frame types
	OpcodeText  = 1
	OpcodeClose = 8
)

// FrameHeader represents a WebSocket frame header
type FrameHeader struct {
	Fin     bool
	Opcode  byte
	Payload []byte
}

// WriteFrame writes a WebSocket frame to the connection
func WriteFrame(w io.Writer, header FrameHeader) error {
	firstByte := header.Opcode
	if header.Fin {
		firstByte |= FinBit
	}

	// Assume payload length is less than 125 for simplicity
	payloadLen := len(header.Payload)

	_, err := w.Write([]byte{firstByte, byte(payloadLen)})
	if err != nil {
		return err
	}

	_, err = w.Write(header.Payload)
	return err
}

// ReadFrame reads a WebSocket frame from the connection
func ReadFrame(r io.Reader) (FrameHeader, error) {
	var header FrameHeader

	// Read first two bytes (FIN + Opcode, and Payload length)
	var b [2]byte
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		return header, err
	}

	header.Fin = b[0]&FinBit != 0
	header.Opcode = b[0] & 0x0F
	payloadLen := int(b[1] & 0x7F)

	// Read the payload data
	header.Payload = make([]byte, payloadLen)
	_, err = io.ReadFull(r, header.Payload)
	return header, err
}
