package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"
)

type CoinMarketResps []struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	Rank             string `json:"rank"`
	PriceUsd         string `json:"price_usd"`
	PriceBtc         string `json:"price_btc"`
	Two4HVolumeUsd   string `json:"24h_volume_usd"`
	MarketCapUsd     string `json:"market_cap_usd"`
	AvailableSupply  string `json:"available_supply"`
	TotalSupply      string `json:"total_supply"`
	MaxSupply        string `json:"max_supply"`
	PercentChange1H  string `json:"percent_change_1h"`
	PercentChange24H string `json:"percent_change_24h"`
	PercentChange7D  string `json:"percent_change_7d"`
	LastUpdated      string `json:"last_updated"`
	PriceCny         string `json:"price_cny"`
	Two4HVolumeCny   string `json:"24h_volume_cny"`
	MarketCapCny     string `json:"market_cap_cny"`
}

type CoinMarket struct {
	Currency  string // CNY, USD, ...
	Interval  int    // in minutes
	Notifiers []NotifierType
}

func (c *CoinMarket) GetName() string {
	return "coinmarketcap"
}

func (c *CoinMarket) GetInterval() time.Duration {
	return time.Minute * time.Duration(c.Interval)
}

func (c *CoinMarket) Fetch() (results []*Item, err error) {
	path := fmt.Sprintf("https://api.coinmarketcap.com/v1/ticker/?convert=%s&limit=10", c.Currency)

	client := &http.Client{}

	req, err := http.NewRequest("GET", path, nil)
	if LogIfErr(err) {
		return
	}

	resp, err := client.Do(req)
	if LogIfErr(err) {
		return
	}

	var cmResps CoinMarketResps
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = decoder.Decode(&cmResps)
	if LogIfErr(err) {
		return
	}

	lines := make([]string, 10)

	for i, m := range cmResps {
		// Bitcoin $13925.90 ¥91564.88 -1.68% -3.68% -10.68% 📉
		var chart string
		if strings.HasPrefix(m.PercentChange1H, "-") {
			chart = "📉"
		} else {
			chart = "📈"
		}

		lines[i] = fmt.Sprintf("%s $%s ¥%s %s%% %s%% %s%% %s",
			m.Name, prettyPriceRound2(m.PriceUsd), prettyPriceRound2(m.PriceCny),
			m.PercentChange1H, m.PercentChange24H, m.PercentChange7D, chart)
	}

	item := &Item{
		Name:       c.GetName(),
		Identifier: c.GetName() + "-summary",
		NotifyText: strings.Join(lines, "\n"),
		ItemFlags:  JustNotify,
		Notifiers:  c.Notifiers,
	}

	results = append(results, item)
	return
}
