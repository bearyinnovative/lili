package house

type HouseDealShenZhen struct {
	*BaseHouseDeal
}

func NewHouseDealShenZhen() *HouseDealShenZhen {
	return &HouseDealShenZhen{
		&BaseHouseDeal{
			cityName:      "深圳",
			cityShortName: "sz",
		},
	}
}
