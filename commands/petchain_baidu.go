package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"
)

type PetChainBaiduResponse struct {
	ErrorNo   string `json:"errorNo"`
	ErrorMsg  string `json:"errorMsg"`
	Timestamp string `json:"timestamp"`
	Data      struct {
		PetsOnSale []struct {
			ID         string `json:"id"`
			PetID      string `json:"petId"`
			BirthType  int    `json:"birthType"`
			Mutation   int    `json:"mutation"`
			Generation int    `json:"generation"`
			RareDegree int    `json:"rareDegree"`
			Desc       string `json:"desc"`
			PetType    int    `json:"petType"`
			Amount     string `json:"amount"`
			BgColor    string `json:"bgColor"`
			PetURL     string `json:"petUrl"`
		} `json:"petsOnSale"`
		TotalCount int  `json:"totalCount"`
		HasData    bool `json:"hasData"`
	} `json:"data"`
}

type PetChainBaidu struct {
	Notifiers []NotifierType
}

func (c *PetChainBaidu) GetName() string {
	return "petchain-baidu"
}

func (c *PetChainBaidu) GetInterval() time.Duration {
	return time.Minute * 2
}

func (c *PetChainBaidu) Fetch() (results []*Item, err error) {
	// https://pet-chain.baidu.com/data/market/queryPetsOnSale (POST https://pet-chain.baidu.com/data/market/queryPetsOnSale)

	jsonData := []byte(`{"pageNo": 1,"pageSize": 100,"petIds": [],"lastAmount": null,"lastRareDegree": null,"appId": 0,"requestId": 1517752179052,"querySortType": "AMOUNT_ASC","tpl": ""}`)
	body := bytes.NewBuffer(jsonData)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "https://pet-chain.baidu.com/data/market/queryPetsOnSale", body)
	if LogIfErr(err) {
		return
	}

	// Headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Origin", "https://pet-chain.baidu.com")
	req.Header.Add("Referer", "https://pet-chain.baidu.com/chain/dogMarket?appId=&tpl=")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "PSTM=1486718302; BIDUPSID=94460EC71ADCA65BB215CC838EE69746; BDUSS=lKbEdsMGZXeXR1YXBGWlpySH50cHRNSjY2MUd1UzRLUmxjVkc5RUhYZnFJVGxhQUFBQUFBJCQAAAAAAAAAAAEAAAAR3sBPaXRwdWJfYmFrAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOqUEVrqlBFaQj; MCITY=340-340%3A; PSINO=7; pgv_pvi=3675588608; pgv_si=s5812241408; BDRCVFR[feWj1Vr5u3D]=mk3SLVN4HKm; BAIDUID=1039E93C7F1FF81859E1557A612CCF70:FG=1; cflag=15%3A3; H_PS_PSSID=25639_1439_24885_21080_17001_22158")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,ja;q=0.6,zh-TW;q=0.5")

	// Fetch Request
	resp, err := client.Do(req)
	if LogIfErr(err) {
		return
	}
	var lbResp PetChainBaiduResponse
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = decoder.Decode(&lbResp)
	if LogIfErr(err) {
		return
	}

	for _, pet := range lbResp.Data.PetsOnSale {
		/*
		   四川龙潮

		   所有者


		   属性
		   体型：皮卡稀有 花纹：峡谷纹 眼睛：小严肃 眼睛色：香苹果 嘴巴：熊出没 肚皮色：白色 身体色：浅蟹灰 花纹色：紫灰
		*/
		ref := fmt.Sprintf("https://pet-chain.baidu.com/chain/detail?from=market&petId=%s&appId=&tpl=", pet.PetID)
		item := &Item{
			Name:       c.GetName(),
			Identifier: c.GetName() + "-" + pet.PetID,
			Desc:       fmt.Sprintf("%s, Generation: %d, Rare: %d, $%s\n%s", pet.Desc, pet.Generation, pet.RareDegree, pet.Amount, ref),
			Ref:        ref,
			Key:        pet.Amount,
			Notifiers:  c.Notifiers,
			ItemFlags:  DoNotCheckTooOld,
		}

		results = append(results, item)
	}

	return
}
