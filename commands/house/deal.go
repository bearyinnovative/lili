package house

import (
	"fmt"
	"log"
	"os"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"
)

var prefetchAll bool

func init() {
	prefetch := os.Getenv("LILI_PREFETCH_ALL_DEALS")
	log.Println("LILI_PREFETCH_ALL_DEALS:", prefetch)
	prefetchAll = prefetch == "1"
}

type HouseDeal struct {
	cityInfo  *CityInfo
	notifiers []NotifierType
}

func NewHouseDeal(name string, notifiers []NotifierType) (*HouseDeal, error) {
	info, err := getCityInfoFromName(name)
	if LogIfErr(err) {
		return nil, err
	}

	return &HouseDeal{
		info, notifiers,
	}, nil
}

func (c *HouseDeal) GetName() string {
	return "house-deal-" + c.cityInfo.Name
}

func (c *HouseDeal) GetInterval() time.Duration {
	return time.Hour * 8
}

func (c *HouseDeal) Fetch() (results []*Item, err error) {
	cityId := c.cityInfo.Id

	stop := false
	offset := 0
	limit := defaultPageCount

	for !stop {
		dealResp, err := fetchDeals(cityId, offset, limit)
		if LogIfErr(err) {
			break
		}

		log.Printf("[%s] fetched %d, has more: %d, total: %d\n",
			c.GetName(),
			len(dealResp.Data.List),
			dealResp.Data.HasMoreData,
			dealResp.Data.TotalCount)

		if dealResp.Errno != 0 {
			log.Printf("[%s] ERROR: %d, %s", c.GetName(), dealResp.Errno, dealResp.Error)
			break
		}

		if len(dealResp.Data.List) == 0 || dealResp.Data.TotalCount == 0 {
			break
		}

		createdCount := 0
		for _, di := range dealResp.Data.List {
			di.CityId = cityId
			created, err := upsertDeal(di)

			if LogIfErr(err) {
				stop = true
				break
			}

			if !created {
				continue
			}

			createdCount += 1

			if len(c.notifiers) == 0 {
				continue
			}

			// start create notify item
			var images []string = nil
			if di.CoverPic != "" {
				images = []string{di.CoverPic}
			}

			// {"title" : "南岭花园 1室1厅 29.24㎡", "price" : 648000, "pricehide" : "6*", "deschide" : "近30天内成交", "unitprice" : 22162, "signdate" : "2017.11.04", "signtimestamp" : 1509788751, "signsource" : "链家成交", "orientation" : "南", "floorstate" : "低楼层/1层", "buildingfinishyear" : 1994, "decoration" : "简装", "buildingtype" : "板楼", "requirelogin" : 0, "fetchedat" : ISODate("2017-11-19T04:29:44.621Z") }
			createdAt := time.Unix(int64(di.SignTimestamp), 0)
			ref := fmt.Sprintf("https://%s.lianjia.com/chengjiao/%s.html", c.cityInfo.Shortname, di.HouseCode)
			item := &Item{
				Name:       c.GetName(),
				Identifier: c.GetName() + "-" + di.HouseCode,
				// 南岭花园 1室1厅 29.24㎡ 南 | 简装 | 低楼层/1层 | 板楼 总价: 648000 单价: 22162 成交时间 2017.11.04
				Desc: fmt.Sprintf("**NEW DEAL** %s %s %s | %s | %s | %s 总价: %.1f万 单价: %.4f万 成交时间: %s [Link](%s)",
					c.cityInfo.Name, di.Title, di.Orientation, di.Decoration, di.FloorState, di.BuildingType, float64(di.Price)/10000.0, float64(di.UnitPrice)/10000.0, di.SignDate, ref),
				Ref:       ref,
				Created:   createdAt,
				Images:    images,
				Notifiers: c.notifiers,
			}

			results = append(results, item)
		}

		if !prefetchAll || createdCount == 0 {
			stop = true
		}

		if !stop {
			offset += limit
		}
	}

	log.Printf("[%s] finished\n", c.GetName())

	return
}
