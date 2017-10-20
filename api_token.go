package seafile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

//获取AuthToken
func (cli *Client) Auth(username, password string) error {
	formData := url.Values{
		"username": {username},
		"password": {password},
	}

	resp, err := http.PostForm(cli.Addr+"/api2/auth-token/", formData)
	if err != nil {
		return fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取错误:%s %s", resp.Status, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s %s", resp.Status, string(b))
	}

	//TODO: 找到真正的返回结构
	respInfo := map[string]string{}
	err = json.Unmarshal(b, &respInfo)
	if err != nil {
		return fmt.Errorf("读取错误:%s", err)
	}

	cli.authToken = respInfo["token"]

	if respInfo["token"] == "" {
		return fmt.Errorf("%s", string(b))
	}

	return nil
}
