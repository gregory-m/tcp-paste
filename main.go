package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nlopes/slack"

	"github.com/gregory-m/tcp-paste/handler"
	"github.com/gregory-m/tcp-paste/server"
)

var runPaste = flag.Bool("paste", true, "Run paste service (save output on local disk)")
var pasteHost = flag.String("paste-host", ":4343", "Host and port for post service")

var runHTTP = flag.Bool("http", true, "Run HTTP service (to server saved output from local disk)")
var httpHost = flag.String("http-host", ":8080", "Host and port for HTTP service")

var storageDir = flag.String("storage", "/tmp", "Storage directory")
var hostName = flag.String("hostname", "localhost:8080", "Hostname to use in links")

var runSlack = flag.Bool("slack", true, "Run Slack service (to post output to slack)")
var slackHost = flag.String("slack-host", ":9393", "Host and port for slack service")
var slackToken = flag.String("slack-token", "", "Slack API token")
var slackChannel = flag.String("slack-chanel", "testa", "Slack API token")

func main() {
	flag.Parse()
	var services []server.Service

	if *runPaste {
		pasteAddr, err := net.ResolveTCPAddr("tcp", *pasteHost)
		if err != nil {
			fmt.Printf("Can't parse paste-host falg: %s\n", err)
			os.Exit(1)
		}

		pasteH := &handler.SaveToDisk{
			HostName:   *hostName,
			StorageDir: *storageDir,
			Prefix:     server.FileServerPrefix,
		}

		pasteS := &server.TCP{
			Addr:    pasteAddr,
			Handler: pasteH,
		}

		services = append(services, pasteS)
	}

	if *runHTTP {
		httpAddr, err := net.ResolveTCPAddr("tcp", *httpHost)
		if err != nil {
			fmt.Printf("Can't parse http-host falg: %s\n", err)
			os.Exit(1)
		}

		httpS := &server.HTTP{
			Host:       httpAddr,
			StorageDir: *storageDir,
			Prefix:     server.FileServerPrefix,
		}

		services = append(services, httpS)
	}

	if *runSlack {
		slackAddr, err := net.ResolveTCPAddr("tcp", *slackHost)
		if err != nil {
			fmt.Printf("Can't parse slack-host falg: %s\n", err)
			os.Exit(1)
		}

		if *slackToken == "" {
			fmt.Print("Can't run slack service without slack token\n")
			os.Exit(1)
		}

		api := slack.New(*slackToken)

		slackH := &handler.Slack{
			Channel: *slackChannel,
			API:     api,
		}

		slcakS := &server.TCP{
			Addr:    slackAddr,
			Handler: slackH,
		}

		services = append(services, slcakS)
	}

	errChan := make(chan error)
	exit := make(chan bool)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	for _, s := range services {
		go func(s server.Service) {
			fmt.Printf("Starting: %s\n", s.Name())
			errChan <- s.Start()
		}(s)
	}

	go func() {
		for range signalChan {
			fmt.Print("Interrupted, stopping services...\n")
			for _, s := range services {
				s.Stop()
			}
			exit <- true
		}
	}()

	<-exit
}
