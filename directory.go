package seafile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

//账户信息
type DirectoryEntry struct {
	Id         string
	Type       string
	Name       string
	Size       int
	Permission string
	Mtime      int
	ParentDir  string `json:"parent_dir"`
}

//列出资料库中指定位置目录的文件和子目录
func (cli *Client) ListDirectoryEntries(libId, path string) ([]DirectoryEntry, error) {
	return cli.ListDirectoryEntriesWithOption(libId, path, nil)
}

//列出资料库中指定位置目录的文件
func (cli *Client) ListDirectoryFileEntries(libId, path string) ([]DirectoryEntry, error) {
	query := url.Values{"t": {"f"}}
	return cli.ListDirectoryEntriesWithOption(libId, path, query)
}

//列出资料库中指定位置目录的子目录
func (cli *Client) ListDirectoryDirectoryEntries(libId, path string) ([]DirectoryEntry, error) {
	query := url.Values{"t": {"d"}}
	return cli.ListDirectoryEntriesWithOption(libId, path, query)
}

//列出资料库中指定位置目录下的所有目录，并递归地获取其子目录下的目录
func (cli *Client) ListDirectoryEntriesRecursive(libId, path string) ([]DirectoryEntry, error) {
	query := url.Values{"t": {"d"}, "recursive": {"1"}}
	return cli.ListDirectoryEntriesWithOption(libId, path, query)
}

//列出资料库指定位置的目录内容
func (cli *Client) ListDirectoryEntriesWithOption(libId, path string, query url.Values) ([]DirectoryEntry, error) {
	if query == nil {
		query = url.Values{}
	}

	if path == "" {
		path = "/"
	}

	query.Set("p", path)

	resp, err := cli.doRequest("GET", "/repos/"+libId+"/dir/?"+query.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(b))

	info := []DirectoryEntry{}
	//err = json.NewDecoder(resp.Body).Decode(&info)
	err = json.Unmarshal(b, &info)
	if err != nil {
		return nil, fmt.Errorf("读取错误:%s %s", resp.Status, err)
	}

	return info, nil
}
