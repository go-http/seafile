package seafile

import (
	"os"
	"testing"
)

func TestDirGet(t *testing.T) {
	client := New(os.Getenv("SEAFILE_HOST"), os.Getenv("SEAFILE_TOKEN"))

	repo, err := client.GetRepoByName("测试")
	if err != nil {
		t.Fatal(err)
	}

	dir, err := repo.GetDir("/文件夹1/子文件夹")
	if err != nil {
		t.Fatalf("获取文件夹失败: %s", err)
	}

	t.Logf("%+v", dir)
}
