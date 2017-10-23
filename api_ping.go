package seafile

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//测试Seafile服务连通性
func (cli *Client) Ping() error {
	resp, err := http.Get(cli.Addr + "/api2/ping")
	if err != nil {
		return fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取错误:%s", err)
	}

	if string(b) == `"pong"` {
		return nil
	}

	return fmt.Errorf("未知返回:%s", err)
}

//自动添加Token后执行请求
func (cli *Client) AuthPing() error {
	resp, err := cli.doRequest("GET", "/auth/ping/", nil, nil)
	if err != nil {
		return fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取错误:%s %s", resp.Status, err)
	}

	if string(b) == `"pong"` {
		return nil
	}

	return fmt.Errorf("未知返回:%s", err)
}
