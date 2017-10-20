package seafile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
func (lib *Library) ListDirectoryEntries(path string) ([]DirectoryEntry, error) {
	return lib.ListDirectoryEntriesWithOption(path, nil)
}

//列出资料库中指定位置目录的文件
func (lib *Library) ListDirectoryFileEntries(path string) ([]DirectoryEntry, error) {
	query := url.Values{"t": {"f"}}
	return lib.ListDirectoryEntriesWithOption(path, query)
}

//列出资料库中指定位置目录的子目录
func (lib *Library) ListDirectoryDirectoryEntries(path string) ([]DirectoryEntry, error) {
	query := url.Values{"t": {"d"}}
	return lib.ListDirectoryEntriesWithOption(path, query)
}

//列出资料库中指定位置目录下的所有目录，并递归地获取其子目录下的目录
func (lib *Library) ListDirectoryEntriesRecursive(path string) ([]DirectoryEntry, error) {
	query := url.Values{"t": {"d"}, "recursive": {"1"}}
	return lib.ListDirectoryEntriesWithOption(path, query)
}

//列出资料库指定位置的目录内容
func (lib *Library) ListDirectoryEntriesWithOption(path string, query url.Values) ([]DirectoryEntry, error) {
	if query == nil {
		query = url.Values{}
	}

	if path == "" {
		path = "/"
	}

	query.Set("p", path)

	resp, err := lib.doRequest("GET", "/dir/?"+query.Encode(), nil, nil)
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

//在资料库创建目录
//  NOTE: 如果指定目录以及存在，会自动创建重命名后的目录，而不会失败
func (lib *Library) CreateDirectory(path string) error {
	query := url.Values{"p": {path}}
	uri := "/dir/?" + query.Encode()

	body := bytes.NewBufferString("operation=mkdir")

	header := http.Header{"Content-Type": {"application/x-www-form-urlencoded"}}

	resp, err := lib.doRequest("POST", uri, header, body)
	if err != nil {
		return fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	return fmt.Errorf("%s %s", resp.Status, string(b))
}

//删除目录
func (lib *Library) RemoveDirectory(dir string) error {
	query := url.Values{"p": {dir}}
	resp, err := lib.doRequest("DELETE", "/dir/?"+query.Encode(), nil, nil)
	if err != nil {
		return fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取错误:%s %s", resp.Status, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("错误:%s %s", resp.Status, string(b))
	}

	return nil
}

//重命名目录
//  NOTE: 如果新目录已经存在，会自动创建重命名后的目录，而不会失败
func (lib *Library) RenameDirectory(path, newname string) error {
	q := url.Values{"p": {path}}

	d := url.Values{
		"operation": {"rename"},
		"newname":   {newname},
	}
	body := bytes.NewBufferString(d.Encode())

	hdr := http.Header{"Content-Type": {"application/x-www-form-urlencoded"}}

	resp, err := lib.doRequest("POST", "/dir/?"+q.Encode(), hdr, body)
	if err != nil {
		return fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(b) == 0 {
		return fmt.Errorf("读取错误: %s", err)
	}

	//FIXME:文档上说返回HTTP 301为成功，实测却是HTTP 200。
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("[%s] %s", resp.Status, string(b))
}
