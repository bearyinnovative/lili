package telegram

import (
	"os"
	"testing"
)

func Test(t *testing.T) {
	token := os.Getenv("TEST_TELEGRAM_TOKEN")
	if token == "" {
		t.Fatal("can't find TEST_TELEGRAM_TOKEN")
	}

	id := os.Getenv("TEST_TELEGRAM_ID")
	if id == "" {
		t.Fatal("can't find TEST_TELEGRAM_ID")
	}

	n := &Notifier{
		token,
		id,
		"markdown",
	}
	err := n.Notify("[hi](http://baidu.com)", []string{
		"https://avatars1.githubusercontent.com/u/1117026?s=40&v=4",
		"https://avatars1.githubusercontent.com/u/1117026?s=40&v=4",
	})
	if err != nil {
		t.Error(err)
	}
}
