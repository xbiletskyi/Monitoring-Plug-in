// Package model contains data structures and handlers for each monitoring scenario. It carries actual business logic of
// the entire plug-in
package model

import (
	"bytes"
	"encoding/binary"
	"errors"
	"monitoring-plug-in/internal/zaplogger"
	"net"
)

// ReadPacket reads the data from Unix socket to a local struct
//
// Parameters:
//   - conn (net.Conn) - socket connection to be read from
//   - packet (interface{}) - monitoring packet interface expected from the socket
func ReadPacket(conn net.Conn, packet interface{}) error {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	packetSize := binary.Size(packet)
	if packetSize <= 0 {
		zaplogger.Logger.Errorf("Invalid packet size")
		return errors.New("invalid packet size")
	}

	buffer := make([]byte, packetSize)
	n, err := conn.Read(buffer)
	if err != nil {
		zaplogger.Logger.Errorf("error reading from socket: %v", err)
		return err
	}
	if n != packetSize {
		zaplogger.Logger.Errorf("incomplete packet. Expected %d bytes, got %d bytes", packetSize, n)
		return errors.New("incomplete packet")
	}

	buf := bytes.NewReader(buffer)
	if err := binary.Read(buf, binary.BigEndian, packet); err != nil {
		zaplogger.Logger.Errorf("error decoding packet: %v", err)
		return err
	}

	zaplogger.Logger.Debugf("Successfully read data: %+v\n", packet)
	return nil
}
