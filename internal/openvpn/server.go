package openvpn

import (
	"bufio"
	"context"
	"fmt"
	"io"
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
	_ = server.connection.Close()
}

func (server *Server) IsConnected() bool {
	return false
}

func (server *Server) RequestStatus(ctx context.Context) (ServerStatus, error) {
	if deadline, ok := ctx.Deadline(); ok {
		_ = server.connection.SetReadDeadline(deadline)
	}

	var status ServerStatus
	status.Clients = make([]ConnectedClient, 0)
	if _, err := server.connection.Write([]byte("status 2\n")); err != nil {
		return ServerStatus{}, err
	}

	for line := range readLines(server.connection) {
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
			if client, err := parseClient(csv[1:]); err != nil {
				fmt.Println(err)
			} else {
				status.Clients = append(status.Clients, client)
			}
		}
	}

	server.lastStatus = status
	server.lastUpdate = time.Now()
	return status, nil
}

func readLines(r io.Reader) <-chan string {
	scanner := bufio.NewScanner(r)
	output := make(chan string)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "END") {
				break
			}
			l := strings.Replace(line, "\r", "", 1)
			output <- l
		}
		close(output)
	}()
	return output
}

func parseClient(csv []string) (ConnectedClient, error) {
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
