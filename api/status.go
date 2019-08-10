package api

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Address string
}

type ConnectedClient struct {
	CommonName         string `json:"name"`
	RealAddress        string `json:"real_address"`
	VirtualAddress     string `json:"virtual_address"`
	VirtualIPv6Address string `json:"virtual_ipv6_address"`
	BytesRX            int    `json:"bytes_rx"`
	BytesTX            int    `json:"bytes_tx"`
	ConnectedSince     int    `json:"connected_since"`
	Username           string `json:"username"`
	ClientID           int    `json:"client_id"`
	PeerID             int    `json:"peer_id"`
}

type ServerStatus struct {
	Version string            `json:"version"`
	Time    int               `json:"time"`
	Clients []ConnectedClient `json:"clients"`
}

func (server *Server) RequestStatus() (ServerStatus, error) {
	d := net.Dialer{Timeout: 5 * time.Second}
	conn, err := d.Dial("tcp", server.Address)
	if err != nil {
		return ServerStatus{}, err
	}
	defer conn.Close()

	// check invite
	if message, err := bufio.NewReader(conn).ReadString('\n'); err != nil {
		return ServerStatus{}, err
	} else {
		if !strings.HasPrefix(message, ">INFO:OpenVPN Management Interface") {
			return ServerStatus{}, fmt.Errorf("server return invalid invite:\n%s", message)
		}
	}

	log.Print("Connected to OpenVPN management")

	var status ServerStatus
	status.Clients = make([]ConnectedClient, 0)
	conn.Write([]byte("status 2\n")) // request status in CSV

	var c = make(chan string)
	go func() {
		buffer := make([]byte, 4096)
		msg := ""
		for {
			n, _ := conn.Read(buffer)
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
