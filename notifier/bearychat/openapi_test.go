package bearychat

import (
	"os"
	"testing"
)

func TestSendOpenAPI(t *testing.T) {
	token := os.Getenv("TEST_BEARYCHAT_TOKEN")
	if token == "" {
		t.Fatal("can't find TEST_BEARYCHAT_TOKEN")
	}

	vid := os.Getenv("TEST_BEARYCHAT_VID")
	if vid == "" {
		t.Fatal("can't find TEST_BEARYCHAT_VID")
	}

	n, err := NewOpenAPINotifier(token, vid)
	if err != nil {
		t.Fatal(err)
	}

	err = n.Notify("hello open api", []string{
		"https://avatars1.githubusercontent.com/u/1117026?s=40&v=4",
		"https://avatars1.githubusercontent.com/u/1117026?s=40&v=4",
	})
	if err != nil {
		t.Error(err)
	}
}
