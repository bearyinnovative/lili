package house

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"
)

const (
	pageCount = 100
)

var cityIdMap map[string]int

func init() {
	cityIdMap = map[string]int{
		"北京":  110000,
		"天津":  120000,
		"上海":  310000,
		"成都":  510100,
		"南京":  320100,
		"杭州":  330100,
		"青岛":  370200,
		"大连":  210200,
		"厦门":  350200,
		"武汉":  420100,
		"深圳":  440300,
		"重庆":  500000,
		"长沙":  430100,
		"西安":  610100,
		"济南":  370101,
		"石家庄": 130100,
		"广州":  440100,
		"东莞":  441900,
		"佛山":  440600,
		"合肥":  340100,
		"烟台":  370600,
		"中山":  442000,
		"珠海":  440400,
		"沈阳":  210100,
		"苏州":  320500,
		"廊坊":  131000,
		"太原":  140100,
		"惠州":  441300,
	}
}

type BaseHouseDeal struct {
	CityName      string
	CityShortName string
	Notifiers     []NotifierType
}

func (c *BaseHouseDeal) GetName() string {
	return "house-deal-" + c.CityName
}

func (c *BaseHouseDeal) GetInterval() time.Duration {
	return time.Hour * 8
}

func (c *BaseHouseDeal) GetNotifiers() []NotifierType {
	return c.Notifiers
}

func (c *BaseHouseDeal) Fetch() (results []*Item, err error) {
	prefetchAll := os.Getenv("LILI_PREFETCH_ALL_DEALS")
	log.Println("LILI_PREFETCH_ALL_DEALS:", prefetchAll)

	cityId, err := getCityIdFromName(c.CityName)
	if LogIfErr(err) {
		return
	}

	stop := false
	offset := 0
	limit := pageCount

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

			if len(c.Notifiers) == 0 {
				continue
			}

			// start create notify item
			var images []string = nil
			if di.CoverPic != "" {
				images = []string{di.CoverPic}
			}

			// {"title" : "南岭花园 1室1厅 29.24㎡", "price" : 648000, "pricehide" : "6*", "deschide" : "近30天内成交", "unitprice" : 22162, "signdate" : "2017.11.04", "signtimestamp" : 1509788751, "signsource" : "链家成交", "orientation" : "南", "floorstate" : "低楼层/1层", "buildingfinishyear" : 1994, "decoration" : "简装", "buildingtype" : "板楼", "requirelogin" : 0, "fetchedat" : ISODate("2017-11-19T04:29:44.621Z") }
			createdAt := time.Unix(int64(di.SignTimestamp), 0)
			ref := fmt.Sprintf("https://%s.lianjia.com/chengjiao/%s.html", c.CityShortName, di.HouseCode)
			item := &Item{
				Name:       c.GetName(),
				Identifier: c.GetName() + "-" + di.HouseCode,
				// 南岭花园 1室1厅 29.24㎡ 南 | 简装 | 低楼层/1层 | 板楼 总价: 648000 单价: 22162 成交时间 2017.11.04
				Desc: fmt.Sprintf("**NEW DEAL** %s %s %s | %s | %s | %s 总价: %.1f万 单价: %.4f万 成交时间: %s [Link](%s)",
					c.CityName, di.Title, di.Orientation, di.Decoration, di.FloorState, di.BuildingType, float64(di.Price)/10000.0, float64(di.UnitPrice)/10000.0, di.SignDate, ref),
				Ref:     ref,
				Created: createdAt,
				Images:  images,
			}

			results = append(results, item)
		}

		if createdCount == 0 {
			stop = true
		}

		if !stop {
			offset += limit
		}
	}

	log.Printf("[%s] finished\n", c.GetName())

	return
}

func getCityIdFromName(name string) (int, error) {
	id, present := cityIdMap[name]
	if !present {
		return -1, errors.New("can't find city id for: " + name)
	}

	return id, nil
}
