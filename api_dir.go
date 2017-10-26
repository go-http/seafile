package seafile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

type Dir struct {
	Id        string `json:"obj_id"` //Dir的ID貌似全是0
	Name      string `json:"obj_name"`
	RepoId    string `json:"repo_id"`
	ParentDir string `json:"parent_dir"`
	Perm      string

	Mtime     time.Time
	Size      int
	FileCount int   `json:"file_count"`
	DirCount  int   `json:"dir_count"`
	repo      *Repo `json:"-"`
}

//文件完整路径
func (dir *Dir) Path() string {
	return filepath.Join(dir.ParentDir, dir.Name)
}

//获取目录
func (repo *Repo) GetDir(path string) (*Dir, error) {
	q := url.Values{"path": {path}}
	resp, err := repo.client.apiGET(repo.Uri() + "/dir/detail/?" + q.Encode())
	if err != nil {
		return nil, fmt.Errorf("请求文件信息失败: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("文件信息错误: %s", resp.Status)
	}

	var detail struct {
		Name      string
		Path      string
		Size      int
		Mtime     time.Time
		RepoId    string `json:"repo_id"`
		FileCount int    `json:"file_count"`
		DirCount  int    `json:"dir_count"`
	}
	err = json.NewDecoder(resp.Body).Decode(&detail)
	if err != nil {
		return nil, fmt.Errorf("解析文件夹信息失败: %s, %s", resp.Status, err)
	}

	fmt.Println(detail.Path)

	dir := Dir{
		//Id: detail.Id,//该接口暂时没有提供该信息
		//Perm: detail.Perm,该接口暂时没有提供该信息
		repo:      repo,
		Name:      detail.Name,
		Mtime:     detail.Mtime,
		RepoId:    detail.RepoId,
		ParentDir: filepath.Dir(path),
		Size:      detail.Size,
		FileCount: detail.FileCount,
		DirCount:  detail.DirCount,
	}

	return &dir, nil
}

//创建文件夹
//如果文件夹已存在，则按照重命名规则创建新文件
func (repo *Repo) Mk(path string) (*Dir, error) {
	q := url.Values{"p": {path}}
	d := url.Values{"operation": {"mkdir"}}
	resp, err := repo.client.apiPOSTForm(repo.Uri()+"/dir/?"+q.Encode(), d)
	if err != nil {
		return nil, fmt.Errorf("请求资料库信息失败: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("[%s]%s", resp.Status, string(b))
	}

	var dir Dir
	err = json.NewDecoder(resp.Body).Decode(&dir)
	if err != nil {
		return nil, fmt.Errorf("解析资料库信息失败: %s, %s", resp.Status, err)
	}

	dir.repo = repo

	return &dir, nil
}

//文件夹内容结构
//TODO:尝试与File、Dir合并
type DirEntry struct {
	Id         string
	Name       string
	Type       string //file或dir
	Mtime      int
	Permission string

	ParentDir string `json:"parent_dir"` //仅在递归获取子目录时有效

	Size                 int    //file entry only
	ModifierName         string `json:"modifier_name"`          //file entry only
	ModifierEmail        string `json:"modifier_email"`         //file entry only
	ModifierContactEmail string `json:"modifier_contact_email"` //file entry only

	IsLocked      bool   `json:"is_locked"`       //file entry only
	LockTime      string `json:"lock_time"`       //file entry only
	LockOwner     string `json:"lock_owner"`      //file entry only
	LockOwnerName string `json:"lock_owner_name"` //file entry only
	LockedByMe    bool   `json:"locked_by_me"`    //file entry only
}

//获取文件夹内容
func (dir *Dir) getEntriesWithOption(t string, recursive bool) ([]DirEntry, error) {
	q := url.Values{
		"t":         {t},
		"p":         {dir.Path()},
		"recursive": {"0"},
	}

	if recursive {
		q.Set("recursive", "1")
	}

	resp, err := dir.repo.client.apiGET(dir.repo.Uri() + "/dir/?" + q.Encode())
	if err != nil {
		return nil, fmt.Errorf("请求资料库信息失败: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("[%s]%s", resp.Status, string(b))
	}

	var entries []DirEntry
	err = json.NewDecoder(resp.Body).Decode(&entries)
	if err != nil {
		return nil, fmt.Errorf("解析资料库信息失败: %s, %s", resp.Status, err)
	}

	return entries, nil
}

//获取文件夹的所有内容
func (dir *Dir) GetEntries() ([]DirEntry, error) {
	return dir.getEntriesWithOption("", false)
}

//获取文件夹下的文件
func (dir *Dir) GetSubFiles() ([]DirEntry, error) {
	return dir.getEntriesWithOption("f", false)
}

//获取文件夹下的子文件夹
func (dir *Dir) GetSubDirs() ([]DirEntry, error) {
	return dir.getEntriesWithOption("d", false)
}

//递归获取文件夹所有的子文件夹
func (dir *Dir) GetSubDirTree() ([]DirEntry, error) {
	return dir.getEntriesWithOption("d", true)
}

//删除文件夹
func (dir *Dir) Delete() error {
	q := url.Values{"p": {dir.Path()}}

	resp, err := dir.repo.client.apiDELETE(dir.repo.Uri() + "/dir/?" + q.Encode())
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
