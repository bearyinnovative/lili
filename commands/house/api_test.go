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
// 	ids := []int{
// 		110000,
// 		310000,
// 		440100,
// 		440300,
// 		120000,
// 		510100,
// 		320100,
// 		330100,
// 		370200,
// 		210200,
// 		350200,
// 		420100,
// 		500000,
// 		430100,
// 		610100,
// 		370101,
// 		130100,
// 		441900,
// 		440600,
// 		340100,
// 		370600,
// 		442000,
// 		440400,
// 		210100,
// 		320500,
// 		131000,
// 		140100,
// 		441300,
// 	}
// 	for _, id := range ids {
// 		t.Log("fetch id:", id)

// 		resp, err := fetchDeals(id, 0, 1)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		t.Log("total:", resp.Data.TotalCount)
// 		for _, d := range resp.Data.List {
// 			t.Log(d)
// 		}
// 	}
// }

// func TestGetHouse(t *testing.T) {
// 	resp, err := fetchHouse(440100, 0, 100)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Logf("got %d %d houses, the last one is: %+v\n", len(resp.Data.List), resp.Data.ReturnCount, resp.Data.List[len(resp.Data.List)-1])
// }

// gz 4381 套
// func TestFetchCommunicaty(t *testing.T) {
// 	resp, err := fetchCommunicates(440100, 0, 20, 0, 20)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Log("total count:", resp.Data.TotalCount)
// }
