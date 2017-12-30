package house

import (
	"strings"
	"testing"
)

func TestAuthKey1(t *testing.T) {
	_, key, err := getAuthKey("https://app.api.lianjia.com/house/chengjiao/search?city_id=440300&limit_offset=0&limit_count=20&request_ts=1510291748")

	if err != nil {
		t.Fatal(err)
	}

	if key != "MjAxNzAzMjRfYW5kcm9pZDozYzVlNWRkMjZlNWY0MTIzYWE5ZjRhOGI5MmM5MjI1MjAzMGJhZjAw" {
		t.Fatal("key error:", key)
	}
}

func TestAuthKey2(t *testing.T) {
	newPath, key, err := getAuthKey("https://app.api.lianjia.com/house/chengjiao/search?city_id=440300&limit_offset=0&limit_count=20")

	if err != nil {
		t.Fatal(err)
	}

	if key == "" {
		t.Fatal("key error:", key)
	}

	if !strings.Contains(newPath, "request_ts=") {
		t.Fatal("path error:", newPath)
	}
}

func TestMakeRequest(t *testing.T) {
	req, err := makeCommonGetRequest("https://app.api.lianjia.com/house/chengjiao/search?city_id=440300&limit_offset=0&limit_count=20")

	if err != nil {
		t.Fatal(err)
	}

	if req == nil {
		t.Fatal("no request")
	}
}

// func TestGetDealList(t *testing.T) {
// 	sendList()
// }

func TestGetHouse(t *testing.T) {
	resp, err := fetchHouse(440100, 0, 100)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("got %d %d houses, the last one is: %+v\n", len(resp.Data.List), resp.Data.ReturnCount, resp.Data.List[len(resp.Data.List)-1])
}
