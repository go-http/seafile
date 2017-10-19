package seafile

import (
	"strings"
)

type Client struct {
	Hostname string
}

const apiPrefix = "/api2"

//新建一个Seafile客户端
func New(hostname string) *Client {
	hostname = strings.TrimSuffix(hostname, "/")
	return &Client{Hostname: hostname + apiPrefix}
}
