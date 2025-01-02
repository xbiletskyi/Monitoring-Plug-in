// Package main serves as an entry point to the program
package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"monitoring-plug-in/internal/model"
	"monitoring-plug-in/internal/socket"
	"monitoring-plug-in/internal/zaplogger"
)

func main() {
	zaplogger.InitLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	socketPaths := []struct {
		path    string
		handler func(conn net.Conn)
	}{
		{"/tmp/monitoring-ue-n11-attach", model.HandleConnectionUeAttach},
	}

	for _, sp := range socketPaths {
		wg.Add(1)
		go socket.StartServer(ctx, sp.path, sp.handler, &wg)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	<-sigChan
	zaplogger.Logger.Info("Received termination signal, shutting down servers...")
	cancel()

	wg.Wait()
	zaplogger.Logger.Info("All servers stopped. Exiting.")
}
