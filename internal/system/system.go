package system

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

// Conn wraps a D-Bus system connection.
type Conn struct {
	conn *dbus.Conn
}

// Connect opens a system bus connection.
func Connect() (*Conn, error) {
	c, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	return &Conn{conn: c}, nil
}

func (c *Conn) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Network provides Wi-Fi/network operations (NetworkManager-backed).
type Network struct{ conn *Conn }

func NewNetwork(c *Conn) *Network { return &Network{conn: c} }

func (n *Network) AvailableNetworks() ([]string, error) {
	if n.conn == nil || n.conn.conn == nil {
		return nil, errors.New("no dbus connection")
	}
	// TODO: call org.freedesktop.NetworkManager.GetDevices and filter wifi.
	return nil, nil
}

// Bluetooth provides Bluetooth operations (BlueZ-backed).
type Bluetooth struct{ conn *Conn }

func NewBluetooth(c *Conn) *Bluetooth { return &Bluetooth{conn: c} }

func (b *Bluetooth) PairedDevices() ([]string, error) {
	if b.conn == nil || b.conn.conn == nil {
		return nil, errors.New("no dbus connection")
	}
	// TODO: enumerate org.bluez devices.
	return nil, nil
}

// Audio provides audio output selection (PipeWire/PulseAudio-backed).
type Audio struct{ conn *Conn }

func NewAudio(c *Conn) *Audio { return &Audio{conn: c} }

func (a *Audio) Outputs() ([]string, error) {
	// TODO: enumerate PipeWire sinks.
	return nil, nil
}
