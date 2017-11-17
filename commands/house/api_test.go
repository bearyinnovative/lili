package house

import (
	"strings"
	"testing"
)

func TestAuthKey1(t *testing.T) {
	_, key, err := getAuthKey("https://app.api.lianjia.com/house/chengjiao/search?city_id=440300&limit_offset=0&limit_count=20&request_ts=1510291748")

	if err != nil {
		t.Error(err)
		return
	}

	if key != "MjAxNzAzMjRfYW5kcm9pZDozYzVlNWRkMjZlNWY0MTIzYWE5ZjRhOGI5MmM5MjI1MjAzMGJhZjAw" {
		t.Error("key error:", key)
		return
	}
}

func TestAuthKey2(t *testing.T) {
	newPath, key, err := getAuthKey("https://app.api.lianjia.com/house/chengjiao/search?city_id=440300&limit_offset=0&limit_count=20")

	if err != nil {
		t.Error(err)
		return
	}

	if key == "" {
		t.Error("key error:", key)
		return
	}

	if !strings.Contains(newPath, "request_ts=") {
		t.Error("path error:", newPath)
	}
}

func TestMakeRequest(t *testing.T) {
	req, err := makeCommonGetRequest("https://app.api.lianjia.com/house/chengjiao/search?city_id=440300&limit_offset=0&limit_count=20")

	if err != nil {
		t.Error(err)
		return
	}

	if req == nil {
		t.Error("no request")
		return
	}
}

// func TestGetDealList(t *testing.T) {
// 	sendList()
// }
