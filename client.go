package seafile

import (
	"fmt"
	"io/ioutil"
	"net/http"
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

//测试Seafile服务连通性
func (cli *Client) Ping() (string, error) {
	resp, err := http.Get(cli.Hostname + "/ping")
	if err != nil {
		return "", fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取错误:%s", err)
	}

	return string(b), nil
}
