package seafile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//账户信息
type Device struct {
	User       string
	DeviceId   string `json:"device_id"`
	DeviceName string `json:"device_name"`

	Key             string
	Platform        string
	PlatformVersion string    `json:"platform_version"`
	ClientVersion   string    `json:"client_version"`
	WipedAt         string    `json:"wiped_at"`
	LastLoginIp     string    `json:"last_login_ip"`
	LastAccessed    time.Time `json:"last_accessed"`
	IsDesktopClient bool      `json:"is_desktop_client"`

	//FIXME:文档中现实还有下面的节点，XD实际上没有
	//SyncedLibraries []struct {
	//	RepoId   string `json:"repo_id"`
	//	SyncTime int    `json:"sync_time"`
	//	RepoName string `json:"repo_name"`
	//} `json:"synced_repos"`
}

//列出资料库中指定位置目录的文件和子目录
func (cli *Client) ListDevices() ([]Device, error) {
	resp, err := cli.doRequest("GET", "/devices/", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	info := []Device{}
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, fmt.Errorf("读取错误:%s %s", resp.Status, err)
	}

	return info, nil
}

//注销设备
func (lib *Library) UnlinkDevice(id, platform string) error {
	d := url.Values{
		"device_id": {id},
		"platform":  {platform},
	}
	body := bytes.NewBufferString(d.Encode())

	hdr := http.Header{"Content-Type": {"application/x-www-form-urlencoded"}}

	resp, err := lib.doRequest("DELETE", "/devices/", hdr, body)
	if err != nil {
		return fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(b) == 0 {
		return fmt.Errorf("读取错误: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("错误:%s %s", resp.Status, string(b))
	}

	return nil
}
