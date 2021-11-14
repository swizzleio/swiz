package network

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Endpoint struct {
	Host        string
	Port        int
	User        string
	FlagPrivate bool
}

func NewEndpointFromHostString(s string) Endpoint {
	endpoint := Endpoint{
		Host: s,
	}
	if parts := strings.Split(endpoint.Host, "@"); len(parts) > 1 {
		endpoint.User = parts[0]
		endpoint.Host = parts[1]
	}
	if parts := strings.Split(endpoint.Host, ":"); len(parts) > 1 {
		endpoint.Host = parts[0]
		endpoint.Port, _ = strconv.Atoi(parts[1])
	}
	return endpoint
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

func (endpoint Endpoint) IsPrivate() bool {
	if endpoint.FlagPrivate {
		return true
	}
	ip := net.ParseIP(endpoint.Host)
	return ip.IsPrivate()
}
