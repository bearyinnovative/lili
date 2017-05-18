package telegram

import (
	"os"
	"strconv"
	"testing"
)

func Test(t *testing.T) {
	token := os.Getenv("TEST_TELEGRAM_TOKEN")
	if token == "" {
		t.Fatal("can't find TEST_TELEGRAM_TOKEN")
	}

	idStr := os.Getenv("TEST_TELEGRAM_ID")
	if idStr == "" {
		t.Fatal("can't find TEST_TELEGRAM_ID")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		t.Fatal(err)
	}

	n := &Notifier{
		token,
		id,
		"Markdown",
	}
	err = n.Notify("[hi](http://baidu.com)", nil)
	if err != nil {
		t.Error(err)
	}
}
