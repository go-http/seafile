package seafile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
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
	resp, err := lib.doRequest("POST", uploadLink+"?ret-json=1", header, body)
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

//更新文件内容
//    fileContentMap的key是文件名，value是文件内容
//当目标文件存在时，会自动重命名上传
func (lib *Library) UpdateFileContent(targetFile string, content []byte) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//填充文件内容
	part, err := writer.CreateFormFile("file", filepath.Base(targetFile))
	if err != nil {
		return fmt.Errorf("创建Multipart错误:%s", err)
	}
	part.Write(content)

	//填充其他字段
	writer.WriteField("target_file", targetFile)

	//FIXME:
	//文档中提到使用relative_path，系统会自动创建不存在的路径
	//但实际上好像没有效果，所以这里暂不支持这个参数
	//writer.WriteField("relative_path", subDir)

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("写Multipart文件错误:%s", err)
	}

	//设置请求Header
	header := http.Header{"Content-Type": {writer.FormDataContentType()}}

	//获取上传地址
	updateLink, err := lib.UpdateLink()
	if err != nil {
		return fmt.Errorf("获取上传地址错误:%s", err)
	}

	//执行上传
	resp, err := lib.doRequest("POST", updateLink, header, body)
	if err != nil {
		return fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取错误:%s %s", resp.Status, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("错误:%s", resp.Status)
	}

	fmt.Println("文件ID", string(b))
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
