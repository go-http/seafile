package seafile

import (
	"encoding/json"
	"fmt"
	"io"
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

//自动添加Token后执行请求
func (cli *Client) doRequest(method, uri string, header http.Header, body io.Reader) (*http.Response, error) {
	//如果外部传进来的不是完整的链接（只是路径），则添上默认的Hostname
	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		uri = cli.Hostname + uri
	}

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求错误:%s", err)
	}

	req.Header.Set("Authorization", "Token "+cli.authToken)
	//如果外部传入Header则设置之
	for k, v := range header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	return http.DefaultClient.Do(req)
}

//自动添加Token后执行请求
func (cli *Client) AuthPing() (string, error) {
	resp, err := cli.doRequest("GET", "/auth/ping", nil, nil)
	if err != nil {
		return "", fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取错误:%s %s", resp.Status, err)
	}

	return string(b), nil
}
