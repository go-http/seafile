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
	"time"
)

type File struct {
	Id        string `json:"obj_id"`
	Name      string `json:"obj_name"`
	Type      string
	Size      int64
	Mtime     time.Time
	RepoId    string `json:"repo_id"`
	IsLocked  bool   `json:"is_locked"`
	ParentDir string `json:"parent_dir"`

	repo *Repo `json:"-"`
}

//文件完整路径
func (file *File) Path() string {
	return filepath.Join(file.ParentDir, file.Name)
}

func (repo *Repo) GetFile(path string) (*File, error) {
	q := url.Values{"p": {path}}
	resp, err := repo.client.apiGET(repo.Uri() + "/file/?" + q.Encode())
	if err != nil {
		return nil, fmt.Errorf("请求文件信息失败: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("文件信息错误: %s", resp.Status)
	}

	var file File
	err = json.NewDecoder(resp.Body).Decode(&file)
	if err != nil {
		return nil, fmt.Errorf("解析文件信息失败: %s, %s", resp.Status, err)
	}

	file.repo = repo

	return &file, nil
}

//检查文件是否存在，如果不存在则创建，返回文件本身
func (repo *Repo) TouchFile(path string) (*File, error) {
	file, err := repo.GetFile(path)
	if err == nil {
		file.repo = repo
		return file, nil
	}

	return repo.CreateFile(path)
}

//创建文件
//如果文件已存在，则按照重命名规则创建新文件
//如果不希望创建重命名的文件，建议使用Repo.Touch方法
func (repo *Repo) CreateFile(path string) (*File, error) {
	q := url.Values{"p": {path}}
	d := url.Values{"operation": {"create"}}
	resp, err := repo.client.apiPOSTForm(repo.Uri()+"/file/?"+q.Encode(), d)
	if err != nil {
		return nil, fmt.Errorf("请求资料库信息失败: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("[%s]%s", resp.Status, string(b))
	}

	var file File
	err = json.NewDecoder(resp.Body).Decode(&file)
	if err != nil {
		return nil, fmt.Errorf("解析资料库信息失败: %s, %s", resp.Status, err)
	}

	file.repo = repo

	return &file, nil
}

//更新文件内容
func (file *File) Update(content []byte) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//填充文件内容
	part, err := writer.CreateFormFile("file", file.Name)
	if err != nil {
		return fmt.Errorf("创建Multipart错误:%s", err)
	}
	part.Write(content)

	//填充其他字段
	writer.WriteField("target_file", file.Path())

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("写Multipart文件错误:%s", err)
	}

	//设置请求Header
	header := http.Header{"Content-Type": {writer.FormDataContentType()}}

	link, err := file.repo.FileUpdateLink()
	if err != nil {
		return fmt.Errorf("获取上传地址错误:%s", err)
	}

	//执行上传
	resp, err := file.repo.client.request("POST", link+"?ret-json=1", header, body)
	if err != nil {
		return fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取错误:%s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[%s] %s", resp.Status, string(b))
	}

	file.Id = string(b)

	return nil
}

//删除文件
func (file *File) Delete() error {
	q := url.Values{"p": {file.Path()}}

	resp, err := file.repo.client.apiDELETE(file.repo.Uri() + "/file/?" + q.Encode())
	if err != nil {
		return fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	b, _ := ioutil.ReadAll(resp.Body)
	return fmt.Errorf("[%s] %s", resp.Status, string(b))
}
