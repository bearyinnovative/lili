package telegram

import "testing"

func Test1(t *testing.T) {
	n := &Notifier{
		"191666900:AAFvM3G5_6mrnOt-aa9GFJ1CUj2ovxd9wI8",
		76344074,
	}
	err := n.Notify("* [hi](http://baidu.com)", nil)
	if err != nil {
		t.Error(err)
	}
}
