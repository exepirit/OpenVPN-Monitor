package openvpn

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Address    string
	connection net.Conn
	lastStatus ServerStatus
	lastUpdate time.Time
}

func (server *Server) Connect() error {
	d := net.Dialer{Timeout: 5 * time.Second}
	var err error
	server.connection, err = d.Dial("tcp", server.Address)
	return err
}

func (server *Server) Close() {
	server.connection.Close()
}

func (server *Server) IsConnected() bool {
	return false
}

func (server *Server) RequestStatus() (ServerStatus, error) {
	var status ServerStatus
	status.Clients = make([]ConnectedClient, 0)
	server.connection.Write([]byte("status 2\n")) // request status in CSV

	var c = make(chan string)
	go func() {
		buffer := make([]byte, 4096)
		msg := ""
		for {
			n, _ := server.connection.Read(buffer)
			msg += string(buffer[:n])
			lines := strings.Split(msg, "\n")
			for _, l := range lines {
				l := strings.Replace(l, "\r", "", 1)
				c <- l
				if strings.HasPrefix(l, "END") {
					return
				}
			}
			msg = lines[len(lines)-1]
		}
	}()

	for {
		line := <-c
		csv := strings.Split(line, ",")

		if len(csv) == 0 {
			continue
		}

		switch csv[0] {
		case "TITLE":
			status.Version = csv[1]
		case "TIME":
			i, _ := strconv.Atoi(csv[2])
			status.Time = i
		case "CLIENT_LIST":
			if client, err := ParseClient(csv[1:]); err != nil {
				fmt.Println(err)
			} else {
				status.Clients = append(status.Clients, client)
			}
		}

		if strings.HasPrefix(line, "END") {
			break
		}
	}

	server.lastStatus = status
	server.lastUpdate = time.Now()
	return status, nil
}

func ParseClient(csv []string) (ConnectedClient, error) {
	if len(csv) != 11 {
		return ConnectedClient{}, fmt.Errorf("invalid response")
	}

	var client ConnectedClient
	client.CommonName = csv[0]
	client.RealAddress = csv[1]
	client.VirtualAddress = csv[2]
	client.VirtualIPv6Address = csv[3]
	i, _ := strconv.Atoi(csv[4])
	client.BytesRX = i
	i, _ = strconv.Atoi(csv[5])
	client.BytesTX = i
	i, _ = strconv.Atoi(csv[7])
	client.ConnectedSince = i
	client.Username = csv[8]

	return client, nil
}
