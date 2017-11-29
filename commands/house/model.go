package house

import (
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
)

type DealResponse struct {
	RequestID string `json:"request_id"`
	Errno     int    `json:"errno"`
	Error     string `json:"error"`
	Data      struct {
		TotalCount  int         `json:"total_count"`
		HasMoreData int         `json:"has_more_data"`
		List        []*DealItem `json:"list"`
	} `json:"data"`
	Cost int `json:"cost"`
}

type DealItem struct {
	BizcircleID         int     `json:"bizcircle_id"`
	CommunityID         int64   `json:"community_id"`
	HouseCode           string  `json:"house_code"`
	Title               string  `json:"title"`
	NewTitle            string  `json:"new_title"`
	KvHouseType         string  `json:"kv_house_type"`
	CoverPic            string  `json:"cover_pic,omitempty"`
	FrameID             string  `json:"frame_id"`
	BlueprintHallNum    int     `json:"blueprint_hall_num"`
	BlueprintBedroomNum int     `json:"blueprint_bedroom_num"`
	Area                float64 `json:"area"`
	Price               int     `json:"price"`
	PriceHide           string  `json:"price_hide"`
	DescHide            string  `json:"desc_hide"`
	UnitPrice           int     `json:"unit_price"`
	SignDate            string  `json:"sign_date"`
	SignTimestamp       int     `json:"sign_timestamp"`
	SignSource          string  `json:"sign_source"`
	Orientation         string  `json:"orientation"`
	FloorState          string  `json:"floor_state"`
	BuildingFinishYear  int     `json:"building_finish_year"`
	Decoration          string  `json:"decoration"`
	BuildingType        string  `json:"building_type,omitempty"`
	RequireLogin        int     `json:"require_login"`

	FetchedAt time.Time
	CityId    int
}

var szHouseNotifiers []NotifierType

func init() {
	szHouseNotifiers = []NotifierType{
		BCChannelNotifier("house_info"),
	}
}
