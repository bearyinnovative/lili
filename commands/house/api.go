package house

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	. "github.com/bearyinnovative/lili/util"
)

const (
	defaultPageCount = 100

	secretKey = "93273ef46a0b880faf4466c48f74878f"
	appID     = "20170324_android"
)

var client *http.Client

func init() {
	// init client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{
		Transport: tr,
	}
}

func fetchDeals(cityId, offset, limit int) (result *DealResponse, err error) {
	path := fmt.Sprintf(
		"https://app.api.lianjia.com/house/chengjiao/search?city_id=%d&limit_offset=%d&limit_count=%d",
		cityId, offset, limit)

	log.Println("fetching", path)

	// Create request
	req, err := makeCommonGetRequest(path)
	if LogIfErr(err) {
		return
	}

	resp, err := client.Do(req)
	if LogIfErr(err) {
		return
	}
	defer resp.Body.Close()

	// // Read Response Body
	// respBody, _ := ioutil.ReadAll(resp.Body)

	// // Display Results
	// fmt.Println("response Status : ", resp.Status)
	// fmt.Println("response Headers : ", resp.Header)
	// fmt.Println("response Body : ", string(respBody))

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&result)
	if LogIfErr(err) {
		return
	}

	return
}

/*
order:
	co21: 价格从低到高
	co22: 价格从高到低
	co41: 单价从低到高
	co12: 面积从大到小
	co32: 最新发布
*/
func fetchHouse(cityId, offset, limit int, comm *CommunityItem) (result *HouseResponse, err error) {
	// https://app.api.lianjia.com/house/ershoufang/searchv4?city_id=440100&priceRequest=&limit_offset=100&communityRequset=&moreRequest=&is_suggestion=0&limit_count=20&sugQueryStr=rs金碧花园第一金碧&comunityIdRequest=c2110343238860955&areaRequest=&is_history=1&schoolRequest=&condition=c2110343238860955rs金碧花园第一金碧&roomRequest=&isFromMap=false&queryStringText=金碧花园第一金碧&request_ts=1514693279
	// https: //app.api.lianjia.com/house/ershoufang/searchv4?comunityIdRequest=c2411100803806&city_id=440300&sugCodition=c2411100803806&is_history=0&limit_offset=0&condition=c2411100803806&queryStringText=%E7%BF%A1%E7%BF%A0%E6%98%8E%E7%8F%A0%E8%8A%B1%E5%9B%AD&isFromMap=false&is_suggestion=1&limit_count=20&request_ts=1514725482

	path := fmt.Sprintf(
		"https://app.api.lianjia.com/house/ershoufang/searchv4"+
			"?city_id=%d&limit_offset=%d&limit_count=%d"+
			"&order=co32&is_history=1"+
			"&sugQueryStr=rs%s&comunityIdRequest=c%s&condition=c%srs%s", // comunityIdRequest is typo as same as the android client
		// "&priceRequest=&tagsText=&communityRequset=&moreRequest=&sugQueryStr=&comunityIdRequest=&areaRequest="+
		// "&schoolRequest=&condition=&roomRequest="+
		// "&is_suggestion=0&is_history=0&isFromMap=false",
		cityId, offset, limit,
		comm.CommunityName, comm.CommunityID, comm.CommunityID, comm.CommunityName)

	log.Println("fetching", path)

	// Create request
	req, err := makeCommonGetRequest(path)
	if LogIfErr(err) {
		return
	}

	resp, err := client.Do(req)
	if LogIfErr(err) {
		return
	}
	defer resp.Body.Close()

	// // Read Response Body
	// respBody, _ := ioutil.ReadAll(resp.Body)

	// // Display Results
	// fmt.Println("response Status : ", resp.Status)
	// fmt.Println("response Headers : ", resp.Header)
	// fmt.Println("response Body : ", string(respBody))

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&result)
	if LogIfErr(err) {
		return
	}

	return
}

// https://app.api.lianjia.com/house/community/search?min_build_year=0&max_build_year=10&city_id=440100&limit_offset=0&limit_count=20&request_ts=1514693664
// limit > 20 会报错
func fetchCommunicates(cityId, offset, limit, minBuildYear, maxBuildYear int) (result *CommunityResponse, err error) {
	path := fmt.Sprintf(
		"https://app.api.lianjia.com/house/community/search"+
			"?city_id=%d&limit_offset=%d&limit_count=%d"+
			"&min_build_year=%d&max_build_year=%d",
		cityId, offset, limit, minBuildYear, maxBuildYear)

	log.Println("fetching", path)

	// Create request
	req, err := makeCommonGetRequest(path)
	if LogIfErr(err) {
		return
	}

	resp, err := client.Do(req)
	if LogIfErr(err) {
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&result)
	if LogIfErr(err) {
		return
	}

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
