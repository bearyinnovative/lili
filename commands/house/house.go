package house

import (
	"errors"
	"fmt"
	"log"
	"strconv"
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
				if hi.CoverPic != "" {
					images = []string{hi.CoverPic}
				}

				ref := fmt.Sprintf("https://%s.lianjia.com/ershoufang/%s.html", c.cityInfo.Shortname, hi.HouseCode)
				item := &Item{
					Name:       c.GetName(),
					Identifier: c.GetName() + "-" + hi.HouseCode,
					// 九龙湖传承别墅 独享私家园林湖景 7室3厅/879.41㎡/南 北/玖珑湖悦源庄 2022.7万
					Desc: fmt.Sprintf("%s %s %s%s\n%s",
						hi.Title, hi.Desc, hi.PriceStr, hi.PriceUnit, ref),
					Ref:        ref,
					Images:     images,
					Key:        convertPrice(hi.Price),
					KeyHistory: hi.historyPriceInStrings(),
					Notifiers:  sub.Notifiers,
					ItemFlags:  DoNotCheckTooOld,
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

func (c *HouseSecondHand) log(format string, v ...interface{}) {
	log.Printf("[%s] "+format, append([]interface{}{c.GetName()}, v...)...)
}
