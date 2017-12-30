package house

import (
	"errors"
	"time"
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

type HouseResponse struct {
	RequestID string `json:"request_id"`
	Errno     int    `json:"errno"`
	Error     string `json:"error"`
	Data      struct {
		TotalCount  int          `json:"total_count"`
		ReturnCount int          `json:"return_count"`
		HasMoreData int          `json:"has_more_data"`
		List        []*HouseItem `json:"list"`
		Subscribe   int          `json:"subscribe"`
		HouseTooFew int          `json:"house_too_few"`
	} `json:"data"`
	Cost int `json:"cost"`
}

type HouseItem struct {
	HouseCode    string `json:"house_code"`
	Title        string `json:"title"`
	Desc         string `json:"desc"`
	PriceStr     string `json:"price_str"`
	PriceUnit    string `json:"price_unit"`
	UnitPriceStr string `json:"unit_price_str"`
	CoverPic     string `json:"cover_pic"`
	CardType     string `json:"card_type"`
	IsFocus      bool   `json:"is_focus"`
	IsVr         bool   `json:"is_vr"`
	BasicList    []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"basic_list"`
	InfoList []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"info_list"`
	CommunityName       string  `json:"community_name"`
	BaiduLa             float64 `json:"baidu_la"`
	BaiduLo             float64 `json:"baidu_lo"`
	BlueprintHallNum    int     `json:"blueprint_hall_num"`
	BlueprintBedroomNum int     `json:"blueprint_bedroom_num"`
	Area                float64 `json:"area"`
	Price               int     `json:"price"`
	UnitPrice           int     `json:"unit_price"`
	ColorTags           []struct {
		Desc  string `json:"desc"`
		Color string `json:"color"`
	} `json:"color_tags,omitempty"`
}

var cityIdMap map[string]int = map[string]int{
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

func getCityIdFromName(name string) (int, error) {
	id, present := cityIdMap[name]
	if !present {
		return -1, errors.New("can't find city id for: " + name)
	}

	return id, nil
}
