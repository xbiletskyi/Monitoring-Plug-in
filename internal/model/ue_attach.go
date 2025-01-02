package model

import (
	"monitoring-plug-in/internal/zaplogger"
	"net"
)

type UeAttachPacket struct {
	UeInfo UeInfo
}

// HandleConnectionUeAttach defines report processing logic
//
// Parameters:
//   - conn (net.Conn) - Unix socket connection to read data from
func HandleConnectionUeAttach(conn net.Conn) {
	packet := UeAttachPacket{}
	ReadPacket(conn, &packet)
	zaplogger.Logger.Infof("New connection attach: %v", packet)
}
