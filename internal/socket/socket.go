// Package socket defines socket server initialization for data receiving
package socket

import (
	"context"
	"monitoring-plug-in/internal/zaplogger"
	"net"
	"os"
	"os/exec"
	"sync"
)

// StartServer initializes a Unix domain socket server and listens for connections
func StartServer(
	ctx context.Context,
	socketPath string,
	handler func(conn net.Conn),
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	if _, err := os.Stat(socketPath); err == nil {
		err := os.Remove(socketPath)
		if err != nil {
			return
		}
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		zaplogger.Logger.Errorf("Error opening socket: %v\n", err)
		return
	}
	defer func() {
		_ = listener.Close()
		_ = os.Remove(socketPath)
	}()

	if err := os.Chmod(socketPath, 0770); err != nil {
		zaplogger.Logger.Fatalf("Failed to change permissions of the socket: %v", err)
	}

	if err := exec.Command("chown", "open5gs:open5gs", socketPath).Run(); err != nil {
		zaplogger.Logger.Fatalf("Failed to change owner of the socket: %v", err)
	}

	zaplogger.Logger.Infof("Listening on %s (UNIX domain socket)...\n", socketPath)

	for {
		connChan := make(chan net.Conn, 1)
		errChan := make(chan error, 1)

		go func() {
			conn, err := listener.Accept()
			if err != nil {
				errChan <- err
			} else {
				connChan <- conn
			}
		}()

		select {
		case <-ctx.Done():
			zaplogger.Logger.Infof("Shutting down server on %s", socketPath)
			return
		case conn := <-connChan:
			zaplogger.Logger.Debug("Connection established...")
			go handler(conn)
		case err := <-errChan:
			zaplogger.Logger.Errorf("Error accepting connection: %v\n", err)
		}
	}
}
