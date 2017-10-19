package seafile

import (
	"encoding/json"
	"fmt"
)

//账户信息
type Account struct {
	Name         string
	Usage        int //已用空间
	Total        int //总空间: -2代表不限制
	Email        string
	Department   string
	Institution  string
	LoginId      string `json:"login_id"`
	ContactEmail string `json:"contact_email"`
}

//自动添加Token后执行请求
func (cli *Client) AccountInfo() (Account, error) {
	resp, err := cli.doRequest("GET", "/account/info/", nil)
	if err != nil {
		return Account{}, fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	var info Account
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return Account{}, fmt.Errorf("读取错误:%s %s", resp.Status, err)
	}

	return info, nil
}
