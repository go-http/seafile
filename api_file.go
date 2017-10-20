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
