package ipc

import (
	"encoding/json"
	"net"
)

// Message is the IPC envelope.
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// Listen creates a unix socket listener.
func Listen(path string) (net.Listener, error) {
	return net.Listen("unix", path)
}

// Dial connects to a unix socket.
func Dial(path string) (net.Conn, error) {
	return net.Dial("unix", path)
}

// Send encodes a message to a connection.
func Send(conn net.Conn, msg Message) error {
	enc := json.NewEncoder(conn)
	return enc.Encode(msg)
}

// Receive decodes a message from a connection.
func Receive(conn net.Conn) (Message, error) {
	var msg Message
	dec := json.NewDecoder(conn)
	err := dec.Decode(&msg)
	return msg, err
}
