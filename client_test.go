package seafile

import (
	"os"
	"testing"
)

func TestPing(t *testing.T) {
	hostname := os.Getenv("SEAFILE_HOST")
	result, err := New(hostname).Ping()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("ping返回", result)
}

func TestAuthToken(t *testing.T) {
	hostname := os.Getenv("SEAFILE_HOST")
	username := os.Getenv("SEAFILE_USER")
	password := os.Getenv("SEAFILE_PASS")

	token, err := New(hostname).AuthToken(username, password)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("获取到Token", token)
}
