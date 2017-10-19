package seafile

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	LibraryTypeMine   = "mine"   //我的资料库类型
	LibraryTypeShared = "shared" //私人共享给我的资料库类型
	LibraryTypeGroup  = "group"  //群组共享给我的资料库类型
	LibraryTypeOrg    = "org"    //公共资料库类型
)

//资料库信息
type Library struct {
	client *Client

	Encrypted bool
	Virtual   bool
	Version   int
	Mtime     int
	Size      int

	Permission     string
	Mtime_relative string `json:"mtime_relative"`
	Owner          string
	Root           string
	Id             string
	Name           string
	Type           string
	HeadCommitId   string `json:"head_commit_id"`
	SizeFormatted  string `json:"size_formatted"`
}

//获取所有可用资料库列表
func (cli *Client) ListAllLibraries() ([]*Library, error) {
	return cli.ListLibrariesByType("")
}

//获取我的资料库列表
func (cli *Client) ListOwnedLibraries() ([]*Library, error) {
	return cli.ListLibrariesByType(LibraryTypeMine)
}

//获取私人共享给我的资料库列表
func (cli *Client) ListSharedLibraries() ([]*Library, error) {
	return cli.ListLibrariesByType(LibraryTypeShared)
}

//获取群组共享的资料库列表
func (cli *Client) ListGroupLibraries() ([]*Library, error) {
	return cli.ListLibrariesByType(LibraryTypeGroup)
}

//获取公共的资料库列表
func (cli *Client) ListOrgLibraries() ([]*Library, error) {
	return cli.ListLibrariesByType(LibraryTypeOrg)
}

//获取指定类型的资料库列表
func (cli *Client) ListLibrariesByType(libType string) ([]*Library, error) {
	path := "/repos/"
	if libType != "" {
		path += "?type=" + libType
	}

	resp, err := cli.doRequest("GET", path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	info := []*Library{}
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, fmt.Errorf("读取错误:%s %s", resp.Status, err)
	}

	for _, lib := range info {
		lib.client = cli
	}

	return info, nil
}

func (lib *Library) doRequest(method, uri string, header http.Header, body io.Reader) (*http.Response, error) {
	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		uri = "/repos/" + lib.Id + uri
	}
	return lib.client.doRequest(method, uri, header, body)
}

//获取资料库的上传地址
func (lib *Library) UploadLink() (string, error) {
	resp, err := lib.doRequest("GET", "/upload-link/", nil, nil)
	if err != nil {
		return "", fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取错误:%s %s", resp.Status, err)
	}

	//返回值是"xxx"格式的，需要去掉头尾的引号
	return string(b[1 : len(b)-1]), nil
}

//获取资料库的更新地址
func (lib *Library) UpdateLink() (string, error) {
	resp, err := lib.doRequest("GET", "/update-link/", nil, nil)
	if err != nil {
		return "", fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取错误:%s %s", resp.Status, err)
	}

	//返回值是"xxx"格式的，需要去掉头尾的引号
	return string(b[1 : len(b)-1]), nil
}
