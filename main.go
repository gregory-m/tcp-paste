package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gregory-m/tcp-paste/handler"
	"github.com/gregory-m/tcp-paste/server"
)

// note, that variables are pointers
var storageDir = flag.String("storage", "/tmp", "Storage directory")
var httpHost = flag.String("http-host", ":8080", "Host and port for HTTP connections")
var tcpHost = flag.String("tcp-host", ":4343", "Host and port for for TCP connections")
var hostname = flag.String("hostname", "localhost:8080", "Hostname to use in links")

func main() {
	flag.Parse()

	tcpAddr, err := net.ResolveTCPAddr("tcp", *tcpHost)
	if err != nil {
		fmt.Printf("Can't parse tcp-host falg: %s", err)
	}

	httpAddr, err := net.ResolveTCPAddr("tcp", *httpHost)
	if err != nil {
		fmt.Printf("Can't parse http-host falg: %s", err)
	}

	saveToDiskH := &handler.SaveToDisk{
		HostName:   *hostname,
		StorageDir: *storageDir,
		Prefix:     server.FileServerPrefix,
	}
	tcpS := server.TCP{
		Addr:    tcpAddr,
		Handler: saveToDiskH,
	}

	httpS := server.HTTP{
		Host:       httpAddr,
		StorageDir: *storageDir,
	}

	errChan := make(chan error)
	signalChan := make(chan os.Signal, 1)
	exit := make(chan bool)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Starting TCP server on: %s\n", *tcpHost)
		errChan <- tcpS.Start()
	}()

	go func() {
		fmt.Printf("Starting HTTP server on: %s\n", *httpHost)
		errChan <- httpS.Start()
	}()

	go func() {
		for range signalChan {
			fmt.Print("Interrupted, stopping services...\n")
			tcpS.Stop()
			exit <- true
		}
	}()

	go func() {
		for range errChan {
			fmt.Printf("Received an error: %s\n", err)
			exit <- true
		}
	}()

	<-exit
}
