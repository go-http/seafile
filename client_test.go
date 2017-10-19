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

	t.Log("Ping返回", result)
}

func TestAuthPing(t *testing.T) {
	hostname := os.Getenv("SEAFILE_HOST")
	username := os.Getenv("SEAFILE_USER")
	password := os.Getenv("SEAFILE_PASS")

	client := New(hostname)
	err := client.Auth(username, password)
	if err != nil {
		t.Fatal(err)
	}

	result, err := client.AuthPing()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Auth ping返回", result)
}

func TestServerInfo(t *testing.T) {
	hostname := os.Getenv("SEAFILE_HOST")
	username := os.Getenv("SEAFILE_USER")
	password := os.Getenv("SEAFILE_PASS")

	client := New(hostname)
	err := client.Auth(username, password)
	if err != nil {
		t.Fatal(err)
	}

	info, err := client.ServerInfo()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", info)
}
