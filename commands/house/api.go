package house

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	. "github.com/bearyinnovative/lili/util"
)

const (
	secretKey = "93273ef46a0b880faf4466c48f74878f"
	appID     = "20170324_android"
)

var client *http.Client

func init() {
	// init client
	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	client = &http.Client{
	// Transport: tr,
	}
}

// just for test
func sendList() {
	// 成交 list (GET https://app.api.lianjia.com/house/chengjiao/search?city_id=440300&limit_offset=0&limit_count=20&request_ts=1510299565)

	// Create request
	req, err := makeCommonGetRequest("https://app.api.lianjia.com/house/chengjiao/search?city_id=440300&limit_offset=0&limit_count=100")
	if LogIfErr(err) {
		return
	}

	fmt.Println(req)

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	// Display Results
	fmt.Println("response Status : ", resp.Status)
	fmt.Println("response Headers : ", resp.Header)
	fmt.Println("response Body : ", string(respBody))
}

func FetchDeals(cityId, offset, limit int) (items []*DealItem, err error) {
	path := fmt.Sprintf(
		"https://app.api.lianjia.com/house/chengjiao/search?city_id=%d&limit_offset=%d&limit_count=%d",
		cityId, offset, limit)

	// Create request
	req, err := makeCommonGetRequest(path)
	if LogIfErr(err) {
		return
	}

	resp, err := client.Do(req)
	if LogIfErr(err) {
		return
	}

	// // Read Response Body
	// respBody, _ := ioutil.ReadAll(resp.Body)

	// // Display Results
	// fmt.Println("response Status : ", resp.Status)
	// fmt.Println("response Headers : ", resp.Header)
	// fmt.Println("response Body : ", string(respBody))

	var result *DealResponse

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = decoder.Decode(&result)
	if LogIfErr(err) {
		return
	}

	items = result.Data.List
	return
}

func makeCommonGetRequest(urlStr string) (*http.Request, error) {
	newUrlStr, authKey, err := getAuthKey(urlStr)
	if LogIfErr(err) {
		return nil, err
	}

	req, err := http.NewRequest("GET", newUrlStr, nil)
	if LogIfErr(err) {
		return nil, err
	}

	// Headers
	req.Header.Add("Connection", "Keep-Alive")
	// req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Lianjia-Version", "7.12.1")
	req.Header.Add("Lianjia-Channel", "Android_Baidu")
	req.Header.Add("Page-Schema", "tradedSearch%2Flist")
	req.Header.Add("User-Agent", "HomeLink7.12.1;Xiaomi MI+6; Android 7.1.1")
	req.Header.Add("Authorization", authKey)
	req.Header.Add("Cookie", "lianjia_udid=865441034412262;lianjia_ssid=70ed25a0-f1dc-4193-b0b0-1969ded8e213;lianjia_uuid=0e5328f3-48b5-4c62-b094-7a0d76f7613e")
	req.Header.Add("Referer", "homepage%3Fcity_id%3D440300")
	req.Header.Add("Host", "app.api.lianjia.com")
	req.Header.Add("Lianjia-Device-Id", "865441034412262")
	req.Header.Add("Lianjia-Im-Version", "2.2.0")

	return req, nil
}

// return new urlStr, auth key, error
func getAuthKey(urlStr string) (string, string, error) {
	u, err := url.Parse(urlStr)
	if LogIfErr(err) {
		return urlStr, "", err
	}

	values := u.Query()

	if values.Get("request_ts") == "" {
		values.Set("request_ts", strconv.FormatInt(time.Now().Unix(), 10))
	}

	u.RawQuery = values.Encode()
	// fmt.Println("new path:", u.String())

	texts := []string{}
	for k, v := range values {
		texts = append(texts, k+"="+v[0])
	}

	sort.Strings(texts)
	// fmt.Println("texts:", texts)
	text := secretKey + strings.Join(texts, "")
	// fmt.Println("before sha1:", text)

	h := sha1.New()
	_, err = h.Write([]byte(text))
	if LogIfErr(err) {
		return urlStr, "", err
	}

	hex := fmt.Sprintf("%x", h.Sum(nil))
	// fmt.Println("hex:", hex)

	authKey := base64.StdEncoding.EncodeToString([]byte(appID + ":" + hex))

	return u.String(), authKey, nil
}
