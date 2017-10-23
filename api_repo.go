package seafile

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Repo struct {
	Id   string `json:"repo_id"`
	Name string `json:"repo_name"`

	OwnerName         string `json:"owner_name"`
	OwnerEmail        string `json:"owner_email"`
	OwnerContactEmail string `json:"owner_contact_email"`
	Encrypted         bool
	Permission        string

	Size      int64
	FileCount int64 `json:"file_count"`

	client *Client `json:"-"`
}

//资料库的资源地址
func (repo *Repo) Uri() string {
	return "/repos/" + repo.Id
}

//根据ID获取资料库信息
func (cli *Client) GetRepo(id string) (*Repo, error) {
	resp, err := cli.apiGET("/repos/" + id + "/")
	if err != nil {
		return nil, fmt.Errorf("请求资料库信息失败: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("资料库信息错误: %s", resp.Status)
	}

	var repo Repo
	err = json.NewDecoder(resp.Body).Decode(&repo)
	if err != nil {
		return nil, fmt.Errorf("解析资料库信息失败: %s, %s", resp.Status, err)
	}

	repo.client = cli

	return &repo, nil
}

//根据name获取资料库信息
func (cli *Client) GetRepoByName(name string) (*Repo, error) {

	var id string
	var err error

	if name == "" {
		id, err = cli.GetDefaultLibraryId()
		if err != nil {
			return nil, fmt.Errorf("获取默认资料资料库ID失败: %s", err)
		}
	} else {
		libraries, err := cli.ListAllLibraries()
		if err != nil {
			return nil, fmt.Errorf("获取资料库列表失败: %s", err)
		}

		for _, library := range libraries {
			if library.Name == name {
				id = library.Id
			}
		}
	}

	if id != "" {
		return cli.GetRepo(id)
	}

	return nil, fmt.Errorf("未找到资料库")
}

//获取资料库的操作地址
func (repo *Repo) getFileServerLink(operation string) (string, error) {
	if operation != "update" && operation != "upload" {
		return "", fmt.Errorf("不支持的操作: %s", operation)
	}

	uri := fmt.Sprintf("%s/%s-link/", repo.Uri(), operation)
	resp, err := repo.client.doRequest("GET", uri, nil, nil)
	if err != nil {
		return "", fmt.Errorf("请求错误:%s", err)
	}
	defer resp.Body.Close()

	var link string
	err = json.NewDecoder(resp.Body).Decode(&link)
	if err != nil {
		return "", fmt.Errorf("解析错误:%s %s", resp.Status, err)
	}

	return link, nil
}

//资料库文件上传地址
func (repo *Repo) FileUploadLink() (string, error) {
	return repo.getFileServerLink("upload")
}

//资料库文件更新地址
func (repo *Repo) FileUpdateLink() (string, error) {
	return repo.getFileServerLink("update")
}
