package openvpn

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
