package handler

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/nlopes/slack"
)

// Slack post connection input to Slack
type Slack struct {
	Channel string
	API     slackAPI
}

type slackAPI interface {
	UploadFile(params slack.FileUploadParameters) (file *slack.File, err error)
}

// Handle implemnts TCPHandler interface interface
func (h *Slack) Handle(conn *net.TCPConn) error {
	defer conn.Close()
	fmt.Printf("%s connected\n", conn.RemoteAddr().String())

	b := new(bytes.Buffer)
	writed, err := io.Copy(b, conn)
	if err != nil {
		return err
	}

	_, err = h.API.UploadFile(slack.FileUploadParameters{
		Title:    "tcp-paste output",
		Filetype: "txt",
		Channels: []string{h.Channel},
		Content:  b.String(),
	})
	if err != nil {
		_, err := conn.Write([]byte(fmt.Sprintf("Can't posted to slack: %s\n", err)))
		return err
	}
	h.API.UploadFile(slack.FileUploadParameters{})
	fmt.Printf("Posted %d bytes for %s at %s\n", writed, conn.RemoteAddr().String(),
		h.Channel)
	_, err = conn.Write([]byte(fmt.Sprintf("Results posted to: %s\n", h.Channel)))
	if err != nil {
		fmt.Printf("can't write to connection: %s\n", err)
		return err
	}
	return nil
}
