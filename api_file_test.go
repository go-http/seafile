package seafile

import (
	"os"
	"testing"
)

func TestFileUpdate(t *testing.T) {
	client := New(os.Getenv("SEAFILE_HOST"), os.Getenv("SEAFILE_TOKEN"))

	repo, err := client.GetRepoByName("测试")
	if err != nil {
		t.Fatal(err)
	}

	file, err := repo.TouchFile("/testdir1/file1.txt")
	if err != nil {
		t.Fatal(err)
	}

	err = file.Update([]byte("我爱北京天安门"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileTouch(t *testing.T) {
	client := New(os.Getenv("SEAFILE_HOST"), os.Getenv("SEAFILE_TOKEN"))

	repo, err := client.GetRepoByName("测试")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("资料库信息")
	t.Logf("%+v", repo)

	file, err := repo.TouchFile("/file123.txt")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", file)
}

func TestDeleteFile(t *testing.T) {
	client := New(os.Getenv("SEAFILE_HOST"), os.Getenv("SEAFILE_TOKEN"))

	repo, err := client.GetRepoByName("测试")
	if err != nil {
		t.Fatalf("获取资料库错误: %s", err)
	}

	file, err := repo.TouchFile("/file_to_be_delete.txt")
	if err != nil {
		t.Fatalf("创建测试文件错误: %s", err)
	}

	t.Logf("已创建准备删除的文件: %+v", file)

	err = file.Delete()
	if err != nil {
		t.Fatal(err)
	}
}
