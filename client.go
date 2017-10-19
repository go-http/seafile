package seafile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	Hostname  string
	authToken string
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

//获取AuthToken
func (cli *Client) AuthToken(username, password string) (string, error) {
	formData := url.Values{
		"username": {username},
		"password": {password},
	}

	resp, err := http.PostForm(cli.Hostname+"/auth-token/", formData)
	if err != nil {
		return "", fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取错误:%s %s", resp.Status, err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s %s", resp.Status, string(b))
	}

	//TODO: 找到真正的返回结构
	respInfo := map[string]string{}
	err = json.Unmarshal(b, &respInfo)
	if err != nil {
		return "", fmt.Errorf("读取错误:%s", err)
	}

	if respInfo["token"] != "" {
		cli.authToken = respInfo["token"]
		return cli.authToken, nil
	} else {
		return "", fmt.Errorf("%s", string(b))
	}
}
