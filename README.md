# OpenVPN-Monitor
Simple OpenVPN client list monitor written on Go

## Installation

Append `management 127.0.0.1 7505` to OpenVPN server configuration file. This enable server management over telnet.

Then run these commands in the shell:

    git clone https://github.com/exepirit/OpenVPN-Monitor
    cd OpenVPN-Monitor
    go build

## Usage

Run `./OpenVPN-Monitor` in directory with configuration and executable file and open in browser [localhost:8080](localhost:8080)
