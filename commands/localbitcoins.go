package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"
)

type LBResponse struct {
	Pagination struct {
		Next string `json:"next"`
	} `json:"pagination"`
	Data struct {
		AdList []struct {
			Data struct {
				Profile struct {
					Username      string    `json:"username"`
					FeedbackScore int       `json:"feedback_score"`
					TradeCount    string    `json:"trade_count"`
					LastOnline    time.Time `json:"last_online"`
					Name          string    `json:"name"`
				} `json:"profile"`
				RequireFeedbackScore       int         `json:"require_feedback_score"`
				HiddenByOpeningHours       bool        `json:"hidden_by_opening_hours"`
				TradeType                  string      `json:"trade_type"`
				AdID                       int         `json:"ad_id"`
				TempPrice                  string      `json:"temp_price"`
				BankName                   string      `json:"bank_name"`
				PaymentWindowMinutes       int         `json:"payment_window_minutes"`
				TrustedRequired            bool        `json:"trusted_required"`
				MinAmount                  string      `json:"min_amount"`
				Visible                    bool        `json:"visible"`
				RequireTrustedByAdvertiser bool        `json:"require_trusted_by_advertiser"`
				TempPriceUsd               string      `json:"temp_price_usd"`
				Lat                        float64     `json:"lat"`
				AgeDaysCoefficientLimit    string      `json:"age_days_coefficient_limit"`
				IsLocalOffice              bool        `json:"is_local_office"`
				FirstTimeLimitBtc          interface{} `json:"first_time_limit_btc"`
				AtmModel                   interface{} `json:"atm_model"`
				City                       string      `json:"city"`
				LocationString             string      `json:"location_string"`
				Countrycode                string      `json:"countrycode"`
				Currency                   string      `json:"currency"`
				LimitToFiatAmounts         string      `json:"limit_to_fiat_amounts"`
				CreatedAt                  time.Time   `json:"created_at"`
				MaxAmount                  string      `json:"max_amount"`
				Lon                        float64     `json:"lon"`
				SmsVerificationRequired    bool        `json:"sms_verification_required"`
				RequireTradeVolume         float64     `json:"require_trade_volume"`
				OnlineProvider             string      `json:"online_provider"`
				MaxAmountAvailable         string      `json:"max_amount_available"`
				Msg                        string      `json:"msg"`
				RequireIdentification      bool        `json:"require_identification"`
				Email                      interface{} `json:"email"`
				VolumeCoefficientBtc       string      `json:"volume_coefficient_btc"`
			} `json:"data"`
			Actions struct {
				PublicView string `json:"public_view"`
			} `json:"actions"`
		} `json:"ad_list"`
		AdCount int `json:"ad_count"`
	} `json:"data"`
}

type BaseLBBuyOnline struct {
	Currency  string // CNY, USD, ...
	Interval  int    // in minutes
	Notifiers []NotifierType
}

func (c *BaseLBBuyOnline) GetName() string {
	return "localbitcoins-buy-online"
}

func (c *BaseLBBuyOnline) GetInterval() time.Duration {
	return time.Minute * time.Duration(c.Interval)
}

func (c *BaseLBBuyOnline) Fetch() (results []*Item, err error) {
	// Request (GET https://localbitcoins.com/buy-bitcoins-online/CNY/.json)
	path := fmt.Sprintf("https://localbitcoins.com/buy-bitcoins-online/%s/.json", c.Currency)

	client := &http.Client{}

	req, err := http.NewRequest("GET", path, nil)
	if LogIfErr(err) {
		return
	}

	resp, err := client.Do(req)
	if LogIfErr(err) {
		return
	}

	var lbResp LBResponse
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = decoder.Decode(&lbResp)
	if LogIfErr(err) {
		return
	}

	for i, ad := range lbResp.Data.AdList {
		data := ad.Data

		// createdAt, err := time.Parse(time.RFC3339, data.CreatedAt)
		if LogIfErr(err) {
			continue
		}

		item := &Item{
			Name:       c.GetName(),
			Identifier: c.GetName() + "-" + strconv.Itoa(data.AdID),
			// vc001 (18; 100%) 124000.00 CNY (500 - 30000) https://localbitcoins.com/ad/644879
			// **124000.00** CNY vc001 (18; 100%) (500 - 30000) https://localbitcoins.com/ad/644879
			Desc:      fmt.Sprintf("[$%s ¥%s %s (%s - %s) %s](%s)", data.TempPriceUsd, data.TempPrice, c.Currency, data.MinAmount, data.MaxAmount, data.Profile.Name, ad.Actions.PublicView),
			Ref:       ad.Actions.PublicView,
			Created:   data.CreatedAt,
			Key:       prettyPriceInWan(data.TempPrice),
			ItemFlags: DoNotCheckTooOld,
		}

		// only notify the lowest price
		if i == 0 {
			item.Notifiers = c.Notifiers
		}

		results = append(results, item)
	}

	return
}
