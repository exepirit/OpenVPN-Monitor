# OpenVPN-Monitor
Simple OpenVPN client list monitor written on Go

## Installation

Append `management 127.0.0.1 7505` to OpenVPN server configuration file. It enable server management over telnet.

```shell script
go get github.com/exepirit/OpenVPN-Monitor/cmd/openvpn-monitor
go install github.com/exepirit/OpenVPN-Monitor/cmd/openvpn-monitor
```

## Usage

1. Run `openvpn-monitor -s remote-server:7505 -b localhost:8080`
2. Open [localhost:8080](http://localhost:8080) in browser.
