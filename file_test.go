package seafile

import (
	"os"
	"testing"
)

func TestGetFile(t *testing.T) {
	client := New(os.Getenv("SEAFILE_HOST"), os.Getenv("SEAFILE_TOKEN"))

	repo, err := client.GetRepoByName(os.Getenv("SEAFILE_REPO"))
	if err != nil {
		t.Fatal(err)
	}

	file, err := repo.GetFile(os.Getenv("SEAFILE_FILE"))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("文件信息:")
	t.Logf("%+v", file)
}
