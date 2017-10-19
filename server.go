package seafile

import (
	"encoding/json"
	"fmt"
)

//账户信息
type Server struct {
	Version  string
	Features []string
}

//获取服务器信息
func (cli *Client) ServerInfo() (Server, error) {
	resp, err := cli.doRequest("GET", "/server-info/", nil, nil)
	if err != nil {
		return Server{}, fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	var info Server
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return Server{}, fmt.Errorf("读取错误:%s %s", resp.Status, err)
	}

	return info, nil
}
