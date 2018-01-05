package house

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"
)

const (
	limit          = defaultPageCount
	communityLimit = 20 // will return error if query with limit more than 20
)

type HouseSubscriber struct {
	Notifiers    []NotifierType
	ShouldNotify func(*HouseItem) bool
}

type HouseSecondHand struct {
	communityOffset    int
	communities        []*CommunityItem
	communityFulfilled bool

	cityInfo    *CityInfo
	subscribers []*HouseSubscriber
}

func NewHouseSecondHand(name string, subscribers []*HouseSubscriber) (*HouseSecondHand, error) {
	info, err := getCityInfoFromName(name)
	if LogIfErr(err) {
		return nil, err
	}

	return &HouseSecondHand{
		0, nil, false, info, subscribers,
	}, nil
}

func (c *HouseSecondHand) GetName() string {
	return "house-secondhand-" + c.cityInfo.Name
}

func (c *HouseSecondHand) GetInterval() time.Duration {
	return time.Second * 60
}

func (c *HouseSecondHand) Fetch() (results []*Item, err error) {
	ci, err := c.loadCommunity()
	if LogIfErr(err) {
		return
	}

	results, err = c.fetchAllHouses(ci)

	return
}

func (c *HouseSecondHand) loadCommunity() (*CommunityItem, error) {
	// if not enough, try load more if could
	if len(c.communities) <= c.communityOffset {
		err := c.loadNextPageCommunityIfNeed()
		LogIfErr(err) // ignore this error, just try fetch houses
	}

	if len(c.communities) == 0 {
		return nil, errors.New("can't find communities")
	}

	defer func() {
		c.communityOffset += 1
	}()

	if c.communityOffset >= len(c.communities) {
		c.log("reset offset from %d", c.communityOffset)
		c.communityOffset = 0
	}

	return c.communities[c.communityOffset], nil
}

func (c *HouseSecondHand) loadNextPageCommunityIfNeed() error {
	// skip if allready fulfilled
	if c.communityFulfilled {
		return nil
	}

	c.log("start loading next page community from %d", len(c.communities))
	resp, err := fetchCommunicates(c.cityInfo.Id, len(c.communities), communityLimit, 0, 20)
	if LogIfErr(err) {
		return err
	}

	if resp.Errno != 0 {
		return fmt.Errorf("code: %d, error: %s", resp.Errno, resp.Error)
	}

	c.communities = append(c.communities, resp.Data.List...)

	if resp.Data.TotalCount <= len(c.communities) {
		c.communityFulfilled = true
		c.log("no more communities total: %d, current: %d", resp.Data.TotalCount, len(c.communities))
	}

	return nil
}

func (c *HouseSecondHand) fetchAllHouses(communityItem *CommunityItem) (results []*Item, err error) {
	offset := 0
	totalChangedCount := 0

	for {
		houseResp, err := fetchHouse(c.cityInfo.Id, offset, limit, communityItem)
		if LogIfErr(err) {
			break
		}

		c.log("fetched %d, has more: %d, total: %d",
			len(houseResp.Data.List),
			houseResp.Data.HasMoreData,
			houseResp.Data.TotalCount)

		if houseResp.Errno != 0 {
			c.log("ERROR: %d, %s", houseResp.Errno, houseResp.Error)
			break
		}

		if len(houseResp.Data.List) == 0 || houseResp.Data.TotalCount == 0 {
			c.log("stop with data count: %d, total count: %d",
				len(houseResp.Data.List), houseResp.Data.TotalCount)
			break
		}

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

			totalChangedCount += 1

			// generate item to notify if need
			for _, sub := range c.subscribers {
				if !sub.ShouldNotify(hi) {
					continue
				}

				// start create notify item
				var images []string = nil
				// 1490190291phpfo1bOJ.png.280x210.jpg is a empty place holder img
				if !strings.Contains(hi.CoverPic, "1490190291phpfo1bOJ.png.280x210.jpg") {
					images = []string{hi.CoverPic}
				}

				// log.Printf("%+v", hi)
				ref := fmt.Sprintf("https://%s.lianjia.com/ershoufang/%s.html", c.cityInfo.Shortname, hi.HouseCode)
				item := &Item{
					Name:       c.GetName(),
					Identifier: c.GetName() + "-" + hi.HouseCode,
					/*
						370万(-2483) 1室1厅/70.27㎡/北/新龙城 回龙观/2006
						https://bj.lianjia.com/ershoufang/101102424514.html
					*/
					Desc: fmt.Sprintf("%s%s(%d) %s %s %s\n%s",
						hi.PriceStr, hi.PriceUnit, hi.UnitPrice-communityItem.ErshoufangAvgUnitPrice,
						// hi.BlueprintHallNum, hi.BlueprintBedroomNum, hi.Area,
						hi.Desc,
						communityItem.BizcircleName, communityItem.BuildingFinishYear,
						ref),
					Ref:        ref,
					Images:     images,
					Key:        convertPrice(hi.Price),
					KeyHistory: hi.historyPriceInStrings(),
					Created:    hi.getCreateTime(),
					Notifiers:  sub.Notifiers,
				}

				results = append(results, item)
			}
		}

		if len(houseResp.Data.List) < limit || offset+limit >= houseResp.Data.TotalCount {
			c.log("stop with data count: %d, total count: %d, offset: %d",
				len(houseResp.Data.List), houseResp.Data.TotalCount, offset)
			break
		}

		offset += limit
	}

	c.log("finished, %d changed", totalChangedCount)

	return
}

func convertPrice(price int) string {
	priceWan := float64(price) / 10000.0
	return strconv.FormatFloat(priceWan, 'f', -1, 64) + "w"
}

func (hi *HouseItem) historyPriceInStrings() []string {
	results := make([]string, len(hi.HistoryPrices))

	for i := 0; i < len(hi.HistoryPrices); i++ {
		results[i] = convertPrice(hi.HistoryPrices[i])
	}

	return results
}

func (hi *HouseItem) getCreateTime() time.Time {
	for _, info := range hi.InfoList {
		if strings.HasPrefix(info.Name, "挂牌") {
			t, err := time.ParseInLocation("2006.01.02", info.Value, time.Local)
			if !LogIfErr(err) {
				return t
			}
		}
	}

	return time.Time{}
}

func (c *HouseSecondHand) log(format string, v ...interface{}) {
	log.Printf("[%s] "+format, append([]interface{}{c.GetName()}, v...)...)
}
