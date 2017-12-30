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

	FetchedAt time.Time `json:"-"`
	CityId    int       `json:"-"`
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

	// db 数据创建时间
	FetchedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	CityId        int   `json:"-"`
	HistoryPrices []int `json:"-"`
}

type CityInfo struct {
	Name      string
	Shortname string
	Id        int
}

var cityInfos []*CityInfo = []*CityInfo{
	&CityInfo{"北京", "bj", 110000},
	&CityInfo{"上海", "sh", 310000},
	&CityInfo{"广州", "gz", 440100},
	&CityInfo{"深圳", "sz", 440300},
	&CityInfo{"天津", "tj", 120000},
	&CityInfo{"成都", "cd", 510100},
	&CityInfo{"南京", "nj", 320100},
	&CityInfo{"杭州", "hz", 330100},
	&CityInfo{"青岛", "qd", 370200},
	&CityInfo{"大连", "dl", 210200},
	&CityInfo{"厦门", "xm", 350200},
	&CityInfo{"武汉", "wh", 420100},
	&CityInfo{"重庆", "cq", 500000},
	&CityInfo{"长沙", "cs", 430100},
	&CityInfo{"西安", "xa", 610100},
	&CityInfo{"济南", "jn", 370101},
	&CityInfo{"石家庄", "sjz", 130100},
	&CityInfo{"东莞", "dg", 441900},
	&CityInfo{"佛山", "fs", 440600},
	&CityInfo{"合肥", "hf", 340100},
	&CityInfo{"烟台", "yt", 370600},
	&CityInfo{"中山", "zs", 442000},
	&CityInfo{"珠海", "zh", 440400},
	&CityInfo{"沈阳", "sy", 210100},
	&CityInfo{"苏州", "s", 320500},
	&CityInfo{"廊坊", "lf", 131000},
	&CityInfo{"太原", "ty", 140100},
	&CityInfo{"惠州", "hui", 441300},
}

func getCityInfoFromName(name string) (*CityInfo, error) {
	for _, info := range cityInfos {
		if info.Name == name {
			return info, nil
		}
	}

	return nil, errors.New("can't find city for: " + name)
}
