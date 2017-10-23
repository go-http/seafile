package seafile

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

//Seafile客户端
type Client struct {
	Addr      string
	authToken string
}

//新建一个Seafile客户端
//三种用法:
//  cli := New(addr)             //未认证的客户端，需要手动调用cli.Auth(user,pass)认证
//  cli := New(addr, token)      //预设置Token的客户端
//  cli := New(addr, user, pass) //带用户信息的客户端，会自动调用Auth以获取Token（忽略错误）
func New(addr string, authParams ...string) *Client {
	client := &Client{
		Addr: strings.TrimSuffix(addr, "/"),
	}

	if len(authParams) == 1 {
		client.authToken = authParams[0]
	} else if len(authParams) == 2 {
		err := client.Auth(authParams[0], authParams[1])
		if err != nil {
			fmt.Printf("用户认证失败: %s", err)
		}
	}

	return client
}

//发起携带Token的Seafile WEB API请求
func (cli *Client) requestApi(apiPrefix, method, uri string, header http.Header, body io.Reader) (*http.Response, error) {
	return cli.request(method, cli.Addr+apiPrefix+uri, header, body)
}

func (cli *Client) request(method, uri string, header http.Header, body io.Reader) (*http.Response, error) {
	//检查Token是否为空
	if cli.authToken == "" {
		return nil, fmt.Errorf("没有合法的Token")
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
