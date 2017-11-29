package bearychat

import "testing"

func TestSendOpenAPI(t *testing.T) {
	n, err := NewOpenAPINotifier("e124ab52a5fe68a6eaa421c9ab3893a3", "=bdwe8zQWz")
	if err != nil {
		t.Error(err)
	}

	err = n.Notify("hello open api", nil)
	if err != nil {
		t.Error(err)
	}
}
