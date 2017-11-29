package bearychat

import "testing"

func TestSendRTM(t *testing.T) {
	n, err := NewRTMNotifier("e124ab52a5fe68a6eaa421c9ab3893a3", "=bdwe8zQWz")
	if err != nil {
		t.Error(err)
	}

	err = n.Notify("hello", []string{
		"https://avatars1.githubusercontent.com/u/1117026?s=40&v=4",
		"https://avatars1.githubusercontent.com/u/1117026?s=40&v=4",
	})

	if err != nil {
		t.Error(err)
	}
}
