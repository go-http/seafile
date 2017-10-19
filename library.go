package seafile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
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

//获取资料库的上传地址
func (cli *Client) LibraryUploadLink(libId string) (string, error) {
	resp, err := cli.doRequest("GET", "/repos/"+libId+"/upload-link/", nil, nil)
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

//上传文件内容
//    fileContentMap的key是文件名，value是文件内容
//当目标文件存在时，会自动重命名上传
func (cli *Client) UploadFileContent(libId, parentDir string, fileContentMap map[string][]byte) error {
	//参数检查
	if !strings.HasSuffix(parentDir, "/") {
		return fmt.Errorf("目录%s必须以/结尾", parentDir)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//填充文件内容
	for filename, content := range fileContentMap {
		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			return fmt.Errorf("创建Multipart错误:%s", err)
		}
		part.Write(content)
	}

	//填充其他字段
	writer.WriteField("parent_dir", parentDir)

	//FIXME:
	//文档中提到使用relative_path，系统会自动创建不存在的路径
	//但实际上好像没有效果，所以这里暂不支持这个参数
	//writer.WriteField("relative_path", subDir)

	err := writer.Close()
	if err != nil {
		return fmt.Errorf("写Multipart文件错误:%s", err)
	}

	//设置请求Header
	header := http.Header{"Content-Type": {writer.FormDataContentType()}}

	//获取上传地址
	uploadLink, err := cli.LibraryUploadLink(libId)
	if err != nil {
		return fmt.Errorf("获取上传地址错误:%s", err)
	}

	//执行上传
	resp, err := cli.doRequest("POST", uploadLink+"?ret-json=1", header, body)
	if err != nil {
		return fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	respInfo := []DirectoryEntry{}
	err = json.NewDecoder(resp.Body).Decode(&respInfo)
	if err != nil {
		return fmt.Errorf("解析错误:%s", err)
	}

	//Debug output
	fmt.Printf("返回%+v\n", respInfo)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("返回%s", resp.Status)
	}

	return nil
}
