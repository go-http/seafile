package seafile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

//上传文件内容
//    fileContentMap的key是文件名，value是文件内容
//当目标文件存在时，会自动重命名上传
func (lib *Library) UploadFileContent(parentDir string, fileContentMap map[string][]byte) error {
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
	uploadLink, err := lib.UploadLink()
	if err != nil {
		return fmt.Errorf("获取上传地址错误:%s", err)
	}

	//执行上传
	resp, err := lib.client.request("POST", uploadLink+"?ret-json=1", header, body)
	if err != nil {
		return fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("[%s] %s", resp.Status, string(b))
	}

	respInfo := []DirectoryEntry{}
	err = json.NewDecoder(resp.Body).Decode(&respInfo)
	if err != nil {
		return fmt.Errorf("解析错误:%s", err)
	}

	return nil
}

//删除文件
func (lib *Library) RemoveFile(file string) error {
	query := url.Values{"p": {file}}
	resp, err := lib.doRequest("DELETE", "/file/?"+query.Encode(), nil, nil)
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

//生成文件下载的链接
//    该链接有效期只有一个小时，过期后无效
//    reuse设置为true时可以不限访问次数，否则访问一次后链接就无效
func (lib *Library) GenerateFileDownloadLink(path string, reuse bool) (string, error) {
	q := url.Values{"p": {path}}
	if reuse {
		q.Set("reuse", "1")
	}

	resp, err := lib.doRequest("GET", "/file/?"+q.Encode(), nil, nil)
	if err != nil {
		return "", fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(b) == 0 {
		return "", fmt.Errorf("读取下载地址错误: %s", err)
	}

	//需要去掉头尾的引号
	return string(b[1 : len(b)-1]), nil
}

//获取文件内容
func (lib *Library) FetchFileContent(path string) ([]byte, error) {
	link, err := lib.GenerateFileDownloadLink(path, false)
	if err != nil {
		return nil, fmt.Errorf("请求下载地址错误:%s", err)
	}

	resp, err := http.Get(link)
	if err != nil {
		return nil, fmt.Errorf("读取文件内容错误: %s", err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

//重命名文件
func (lib *Library) RenameFile(path, newname string) error {
	q := url.Values{"p": {path}}

	d := url.Values{
		"operation": {"rename"},
		"newname":   {newname},
	}
	body := bytes.NewBufferString(d.Encode())

	hdr := http.Header{"Content-Type": {"application/x-www-form-urlencoded"}}

	resp, err := lib.doRequest("POST", "/file/?"+q.Encode(), hdr, body)
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

//复制文件到另一个资料库的指定目录
//Note:
//  目标目录必须存在
//  目标目录下如果有同名文件，新文件会自动重命名
func (lib *Library) CopyFileToLibrary(path, dstLibId, dstLibPath string) error {
	q := url.Values{"p": {path}}

	d := url.Values{
		"operation": {"copy"},
		"dst_repo":  {dstLibId},
		"dst_dir":   {dstLibPath},
	}
	body := bytes.NewBufferString(d.Encode())

	hdr := http.Header{"Content-Type": {"application/x-www-form-urlencoded"}}

	resp, err := lib.doRequest("POST", "/file/?"+q.Encode(), hdr, body)
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

//复制文件到另一个资料库的指定目录，目标目录必须存在
func (lib *Library) MoveFileToLibrary(path, dstLibId, dstLibPath string) error {
	q := url.Values{"p": {path}}

	d := url.Values{
		"operation": {"move"},
		"dst_repo":  {dstLibId},
		"dst_dir":   {dstLibPath},
	}
	body := bytes.NewBufferString(d.Encode())

	hdr := http.Header{"Content-Type": {"application/x-www-form-urlencoded"}}

	resp, err := lib.doRequest("POST", "/file/?"+q.Encode(), hdr, body)
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
