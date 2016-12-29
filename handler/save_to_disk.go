package handler

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/url"
	"os"
	"path"
	"time"
)

// SaveToDisk save connection input to disk
type SaveToDisk struct {
	HostName   string
	StorageDir string
	Prefix     string
}

// Handle implemnts TCPHandler interface interface
func (h *SaveToDisk) Handle(conn *net.TCPConn) error {
	defer conn.Close()

	fmt.Printf("%s connected\n", conn.RemoteAddr().String())
	defer conn.Close()
	fileName := randomString(10)

	file, err := os.Create(path.Join(h.StorageDir, fileName))
	defer file.Close()
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("Sorry can't create file: %s", err)))
		return err
	}

	writed, err := io.Copy(file, conn)
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("Sorry write to file: %s", err)))
		return err
	}

	u := &url.URL{
		Scheme: "http",
		Host:   h.HostName,
		Path:   h.Prefix,
	}
	u.Path = path.Join(u.Path, fileName)

	fmt.Printf("Saved %d bytes for %s at %s\n", writed, conn.RemoteAddr().String(),
		u.String())

	_, err = conn.Write([]byte(fmt.Sprintf("Results saved at: %s\n", u.String())))
	if err != nil {
		fmt.Printf("can't write to connection: %s\n", err)
	}

	return nil
}

func randomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UTC().UnixNano())

	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
