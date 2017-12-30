package house

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"
)

const (
	limit = defaultPageCount
)

type HouseSubscriber struct {
	Notifiers    []NotifierType
	ShouldNotify func(*HouseItem) bool
}

type HouseSecondHand struct {
	offset      int
	stopped     bool
	cityInfo    *CityInfo
	subscribers []*HouseSubscriber
}

func NewHouseSecondHand(name string, subscribers []*HouseSubscriber) (*HouseSecondHand, error) {
	info, err := getCityInfoFromName(name)
	if LogIfErr(err) {
		return nil, err
	}

	return &HouseSecondHand{
		0, false, info, subscribers,
	}, nil
}

func (c *HouseSecondHand) GetName() string {
	return "house-secondhand-" + c.cityInfo.Name
}

func (c *HouseSecondHand) GetInterval() time.Duration {
	return time.Second*60 + time.Duration(rand.Intn(60))*time.Second
	// return time.Second * 5
}

func (c *HouseSecondHand) Fetch() (results []*Item, err error) {
	if c.stopped {
		c.log("skip for stopped")
		return nil, nil
	}

	houseResp, err := fetchHouse(c.cityInfo.Id, c.offset, limit)
	if LogIfErr(err) {
		return
	}

	c.log("fetched %d, has more: %d, total: %d",
		len(houseResp.Data.List),
		houseResp.Data.HasMoreData,
		houseResp.Data.TotalCount)

	if houseResp.Errno != 0 {
		c.log("ERROR: %d, %s", houseResp.Errno, houseResp.Error)

		// 20003, limit_offset is invalid.
		if houseResp.Errno == 20003 {
			c.log("20003 offset invalid, reset")
			c.offset = 0
		}

		return
	}

	if houseResp.Data.TotalCount == 0 {
		c.log("total count == 0, stopped")
		c.stopped = true
		return
	}

	// no more results? reset offset
	if len(houseResp.Data.List) < limit {
		c.log("reset offset with data(%d) < limit(%d)", len(houseResp.Data.List), limit)
		c.offset = 0
	} else if c.offset > houseResp.Data.TotalCount {
		c.log("reset offset with offset(%d) > total count(%d)", c.offset, houseResp.Data.TotalCount)
		c.offset = 0
	} else {
		c.offset += limit
	}

	changedCount := 0
	for _, hi := range houseResp.Data.List {
		hi.CityId = c.cityInfo.Id

		changed := false
		changed, err = upsertHouse(hi)

		if LogIfErr(err) {
			continue
		}

		if !changed {
			continue
		}

		changedCount += 1

		// generate item to notify if need
		for _, sub := range c.subscribers {
			if !sub.ShouldNotify(hi) {
				continue
			}

			// start create notify item
			var images []string = nil
			if hi.CoverPic != "" {
				images = []string{hi.CoverPic}
			}

			ref := fmt.Sprintf("https://%s.lianjia.com/ershoufang/%s.html", c.cityInfo.Shortname, hi.HouseCode)
			item := &Item{
				Name:       c.GetName(),
				Identifier: c.GetName() + "-" + hi.HouseCode,
				// 南岭花园 1室1厅 29.24㎡ 南 | 简装 | 低楼层/1层 | 板楼 总价: 648000 单价: 22162 成交时间 2017.11.04
				// 九龙湖传承别墅 独享私家园林湖景 7室3厅/879.41㎡/南 北/玖珑湖悦源庄 2022.7万
				Desc: fmt.Sprintf("%s %s %s%s\n%s",
					hi.Title, hi.Desc, hi.PriceStr, hi.PriceUnit, ref),
				Ref:        ref,
				Images:     images,
				Key:        hi.PriceStr,
				KeyHistory: hi.historyPriceInStrings(),
				Notifiers:  sub.Notifiers,
				ItemFlags:  DoNotCheckTooOld,
			}

			results = append(results, item)
		}
	}

	c.log("finished, %d changed", changedCount)

	return
}

func (hi *HouseItem) historyPriceInStrings() []string {
	results := make([]string, len(hi.HistoryPrices))

	for i := 0; i < len(hi.HistoryPrices); i++ {
		results[i] = strconv.Itoa(hi.HistoryPrices[i])
	}

	return results
}

func (c *HouseSecondHand) log(format string, v ...interface{}) {
	log.Printf("[%s] "+format, append([]interface{}{c.GetName()}, v...)...)
}
