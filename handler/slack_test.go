package handler

import (
	"net"
	"testing"
	"time"

	"github.com/gregory-m/tcp-paste/server"
	"github.com/nlopes/slack"
)

var testOUT = make(chan slack.FileUploadParameters, 1)

func TestUpload(t *testing.T) {
	testChannelName := "Testa"
	testInput := "input bla-bla-bla"

	ts := newSlackTestSerevr(testChannelName)
	go ts.Start()
	time.Sleep(2 * time.Millisecond)

	c, err := net.Dial("tcp", "127.0.0.1:9393")
	if err != nil {
		t.Fatalf("Can't dial to server: %s", err)
	}

	_, err = c.Write([]byte(testInput))
	if err != nil {
		t.Fatalf("Can't write: %s", err)
	}
	c.(*net.TCPConn).CloseWrite()

	up := <-testOUT

	if len(up.Channels) != 1 {
		t.Fatalf("Uploaded to %d channels", len(up.Channels))
	}

	if up.Channels[0] != testChannelName {
		t.Errorf("Uploaded to wrong channel.\nExpected: %q\nGot:      %q\n",
			testChannelName, up.Channels[0])
	}

	if up.Content != testInput {
		t.Errorf("Content don't match.\nExpected: %q\nGot:      %q\n",
			testInput, up.Content)
	}
}

type slackMock struct{}

func (s *slackMock) UploadFile(params slack.FileUploadParameters) (file *slack.File, err error) {
	testOUT <- params
	return &slack.File{}, nil
}

func newSlackTestHandler(channel string) server.TCPHandler {
	return &Slack{
		Channel: channel,
		API:     &slackMock{},
	}
}

func newSlackTestSerevr(channel string) *server.TCP {
	h := newSlackTestHandler(channel)
	a, err := net.ResolveTCPAddr("tcp", ":9393")
	if err != nil {
		panic(err)
	}

	s := &server.TCP{
		Handler: h,
		Addr:    a,
	}

	return s
}
