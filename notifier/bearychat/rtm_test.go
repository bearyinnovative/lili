package bearychat

import "testing"

func TestSendRTM(t *testing.T) {
	n, err := NewRTMNotifier("e124ab52a5fe68a6eaa421c9ab3893a3", "=bdwe8zQWz")
	if err != nil {
		t.Error(err)
	}

	err = n.Notify("hello", nil)
	if err != nil {
		t.Error(err)
	}
}
