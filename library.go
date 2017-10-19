package seafile

import (
	"encoding/json"
	"fmt"
)

const (
	LibraryTypeMine   = "mine"   //我的资料库类型
	LibraryTypeShared = "shared" //私人共享给我的资料库类型
	LibraryTypeGroup  = "group"  //群组共享给我的资料库类型
	LibraryTypeOrg    = "org"    //公共资料库类型
)

//账户信息
type Library struct {
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
func (cli *Client) ListAllLibraries() ([]Library, error) {
	return cli.ListLibrariesByType("")
}

//获取我的资料库列表
func (cli *Client) ListOwnedLibraries() ([]Library, error) {
	return cli.ListLibrariesByType(LibraryTypeMine)
}

//获取私人共享给我的资料库列表
func (cli *Client) ListSharedLibraries() ([]Library, error) {
	return cli.ListLibrariesByType(LibraryTypeShared)
}

//获取群组共享的资料库列表
func (cli *Client) ListGroupLibraries() ([]Library, error) {
	return cli.ListLibrariesByType(LibraryTypeGroup)
}

//获取公共的资料库列表
func (cli *Client) ListOrgLibraries() ([]Library, error) {
	return cli.ListLibrariesByType(LibraryTypeOrg)
}

//获取指定类型的资料库列表
func (cli *Client) ListLibrariesByType(libType string) ([]Library, error) {
	path := "/repos/"
	if libType != "" {
		path += "?type=" + libType
	}

	resp, err := cli.doRequest("GET", path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	info := []Library{}
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, fmt.Errorf("读取错误:%s %s", resp.Status, err)
	}

	return info, nil
}
