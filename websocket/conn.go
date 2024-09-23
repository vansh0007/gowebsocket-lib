package websocket

import (
	"bufio"
	"errors"
	"io"
	"net"
)

// Conn represents the WebSocket connection
type Conn struct {
	conn net.Conn
}

// NewConn creates a new WebSocket connection
func NewConn(conn net.Conn) *Conn {
	return &Conn{conn: conn}
}

// ReadMessage reads a WebSocket message from the connection.
func (c *Conn) ReadMessage() (string, []byte, error) {
	reader := bufio.NewReader(c.conn)

	// Read the first byte (FIN, opcode)
	finOpcode, err := reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			return "", nil, errors.New("connection closed by client")
		}
		return "", nil, err
	}
	fin := (finOpcode & 0x80) != 0 // FIN bit
	opcode := finOpcode & 0x0F     // Opcode

	// Handle close frames (opcode 8)
	if opcode == 8 {
		return "", nil, errors.New("client closed connection")
	}

	// Handle text frames (opcode 1) and binary frames (opcode 2)
	var messageType string
	if opcode == 1 {
		messageType = "text"
	} else if opcode == 2 {
		messageType = "binary"
	} else {
		return "", nil, errors.New("unsupported frame type")
	}

	// Read the second byte (Mask, payload length)
	maskPayloadLen, err := reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			return "", nil, errors.New("connection closed by client")
		}
		return "", nil, err
	}
	mask := (maskPayloadLen & 0x80) != 0 // Mask bit
	payloadLen := int(maskPayloadLen & 0x7F)

	// Handle extended payload lengths (126, 127)
	if payloadLen == 126 {
		lenBytes := make([]byte, 2)
		_, err := io.ReadFull(reader, lenBytes)
		if err != nil {
			return "", nil, err
		}
		payloadLen = int(lenBytes[0])<<8 | int(lenBytes[1])
	}

	if payloadLen == 127 {
		lenBytes := make([]byte, 8)
		_, err := io.ReadFull(reader, lenBytes)
		if err != nil {
			return "", nil, err
		}
		return "", nil, errors.New("large payloads not supported")
	}

	// Read the masking key if present
	var maskKey []byte
	if mask {
		maskKey = make([]byte, 4)
		_, err := io.ReadFull(reader, maskKey)
		if err != nil {
			return "", nil, err
		}
	}

	// Read the payload
	payload := make([]byte, payloadLen)
	_, err = io.ReadFull(reader, payload)
	if err != nil {
		if err == io.EOF {
			return "", nil, errors.New("connection closed by client")
		}
		return "", nil, err
	}

	// Unmask the payload if masked
	if mask {
		for i := 0; i < payloadLen; i++ {
			payload[i] ^= maskKey[i%4]
		}
	}

	// Ensure that the message is valid UTF-8 if it is a text frame
	if opcode == 1 && !fin {
		return "", nil, errors.New("fragmented messages are not supported")
	}

	return messageType, payload, nil
}

// WriteMessage writes a message to the WebSocket connection.
func (c *Conn) WriteMessage(msg []byte) error {
	var opcode byte = 1 // Default to text frame

	// Determine if the message is binary or text
	if isBinaryMessage(msg) {
		opcode = 2 // Set to binary frame
	}

	if len(msg) > 125 {
		return errors.New("message too large")
	}

	// Prepare the frame header
	frame := []byte{
		0x80 | opcode,  // FIN = 1, Opcode (text or binary)
		byte(len(msg)), // Payload length
	}

	// Append the message to the frame
	frame = append(frame, msg...)

	// Write the frame to the connection
	_, err := c.conn.Write(frame)
	return err
}

// isBinaryMessage determines if the message should be treated as binary.
func isBinaryMessage(msg []byte) bool {
	// Check for non-printable characters
	for _, b := range msg {
		if b < 32 || b > 126 { // ASCII range for printable characters
			return true // Treat as binary if non-printable
		}
	}
	return false // Otherwise, treat as text
}

// Close closes the WebSocket connection.
func (c *Conn) Close() error {
	return c.conn.Close()
}
