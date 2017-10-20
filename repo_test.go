package seafile

import (
	"os"
	"testing"
)

func TestGetRepoByName(t *testing.T) {
	client := New(os.Getenv("SEAFILE_HOST"), os.Getenv("SEAFILE_TOKEN"))

	repo, err := client.GetRepoByName(os.Getenv("SEAFILE_REPO"))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("资料库信息")
	t.Logf("%+v", repo)
}
