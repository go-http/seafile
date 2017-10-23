package seafile

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
		return nil, fmt.Errorf("资料库信息错误: %s", resp.Status)
	}

	var file File
	err = json.NewDecoder(resp.Body).Decode(&file)
	if err != nil {
		return nil, fmt.Errorf("解析资料库信息失败: %s, %s", resp.Status, err)
	}

	file.repo = repo

	return &file, nil
}
