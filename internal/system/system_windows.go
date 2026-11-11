//go:build windows

package system

import "errors"

// Conn stub for Windows (no D-Bus).
type Conn struct{}

func Connect() (*Conn, error) { return &Conn{}, nil }

func (c *Conn) Close() error { return nil }

type Network struct{ conn *Conn }

func NewNetwork(c *Conn) *Network { return &Network{conn: c} }

func (n *Network) AvailableNetworks() ([]string, error) {
	return nil, errors.New("not supported on Windows")
}

type Bluetooth struct{ conn *Conn }

func NewBluetooth(c *Conn) *Bluetooth { return &Bluetooth{conn: c} }

func (b *Bluetooth) PairedDevices() ([]string, error) {
	return nil, errors.New("not supported on Windows")
}

type Audio struct{ conn *Conn }

func NewAudio(c *Conn) *Audio { return &Audio{conn: c} }

func (a *Audio) Outputs() ([]string, error) {
	return nil, errors.New("not supported on Windows")
}
