package seafile

import (
	"os"
	"testing"
)

func TestPing(t *testing.T) {
	hostname := os.Getenv("SEAFILE_HOST")
	result, err := New(hostname).Ping()
	if err != nil {
		t.Error(err)
	}

	t.Log(result)
}
