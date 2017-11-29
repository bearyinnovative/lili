package house

type HouseDealBeiJing struct {
	*BaseHouseDeal
}

func NewHouseDealBeiJing() *HouseDealBeiJing {
	return &HouseDealBeiJing{
		&BaseHouseDeal{
			"北京", "bj", nil,
		},
	}
}

type HouseDealShangHai struct {
	*BaseHouseDeal
}

func NewHouseDealShangHai() *HouseDealShangHai {
	return &HouseDealShangHai{
		&BaseHouseDeal{
			"上海", "sh", nil,
		},
	}
}

type HouseDealGuangZhou struct {
	*BaseHouseDeal
}

func NewHouseDealGuangZhou() *HouseDealGuangZhou {
	return &HouseDealGuangZhou{
		&BaseHouseDeal{
			"广州", "gz", nil,
		},
	}
}

type HouseDealShenZhen struct {
	*BaseHouseDeal
}

func NewHouseDealShenZhen() *HouseDealShenZhen {
	return &HouseDealShenZhen{
		&BaseHouseDeal{
			"深圳", "sz", szHouseNotifiers,
		},
	}
}

type HouseDealTJ struct {
	*BaseHouseDeal
}

func NewHouseDealTJ() *HouseDealTJ {
	return &HouseDealTJ{
		&BaseHouseDeal{
			"天津", "tj", nil,
		},
	}
}

type HouseDealCD struct {
	*BaseHouseDeal
}

func NewHouseDealCD() *HouseDealCD {
	return &HouseDealCD{
		&BaseHouseDeal{
			"成都", "cd", nil,
		},
	}
}

type HouseDealNJ struct {
	*BaseHouseDeal
}

func NewHouseDealNJ() *HouseDealNJ {
	return &HouseDealNJ{
		&BaseHouseDeal{
			"南京", "nj", nil,
		},
	}
}

type HouseDealHZ struct {
	*BaseHouseDeal
}

func NewHouseDealHZ() *HouseDealHZ {
	return &HouseDealHZ{
		&BaseHouseDeal{
			"杭州", "hz", nil,
		},
	}
}

type HouseDealQD struct {
	*BaseHouseDeal
}

func NewHouseDealQD() *HouseDealQD {
	return &HouseDealQD{
		&BaseHouseDeal{
			"青岛", "qd", nil,
		},
	}
}

type HouseDealDL struct {
	*BaseHouseDeal
}

func NewHouseDealDL() *HouseDealDL {
	return &HouseDealDL{
		&BaseHouseDeal{
			"大连", "dl", nil,
		},
	}
}

type HouseDealXM struct {
	*BaseHouseDeal
}

func NewHouseDealXM() *HouseDealXM {
	return &HouseDealXM{
		&BaseHouseDeal{
			"厦门", "xm", nil,
		},
	}
}

type HouseDealWH struct {
	*BaseHouseDeal
}

func NewHouseDealWH() *HouseDealWH {
	return &HouseDealWH{
		&BaseHouseDeal{
			"武汉", "wh", nil,
		},
	}
}

type HouseDealCQ struct {
	*BaseHouseDeal
}

func NewHouseDealCQ() *HouseDealCQ {
	return &HouseDealCQ{
		&BaseHouseDeal{
			"重庆", "cq", nil,
		},
	}
}

type HouseDealCS struct {
	*BaseHouseDeal
}

func NewHouseDealCS() *HouseDealCS {
	return &HouseDealCS{
		&BaseHouseDeal{
			"长沙", "cs", nil,
		},
	}
}

type HouseDealXA struct {
	*BaseHouseDeal
}

func NewHouseDealXA() *HouseDealXA {
	return &HouseDealXA{
		&BaseHouseDeal{
			"西安", "xa", nil,
		},
	}
}

type HouseDealJN struct {
	*BaseHouseDeal
}

func NewHouseDealJN() *HouseDealJN {
	return &HouseDealJN{
		&BaseHouseDeal{
			"济南", "jn", nil,
		},
	}
}

type HouseDealSJZ struct {
	*BaseHouseDeal
}

func NewHouseDealSJZ() *HouseDealSJZ {
	return &HouseDealSJZ{
		&BaseHouseDeal{
			"石家庄", "sjz", nil,
		},
	}
}

type HouseDealDG struct {
	*BaseHouseDeal
}

func NewHouseDealDG() *HouseDealDG {
	return &HouseDealDG{
		&BaseHouseDeal{
			"东莞", "dg", nil,
		},
	}
}

type HouseDealFS struct {
	*BaseHouseDeal
}

func NewHouseDealFS() *HouseDealFS {
	return &HouseDealFS{
		&BaseHouseDeal{
			"佛山", "fs", nil,
		},
	}
}

type HouseDealHF struct {
	*BaseHouseDeal
}

func NewHouseDealHF() *HouseDealHF {
	return &HouseDealHF{
		&BaseHouseDeal{
			"合肥", "hf", nil,
		},
	}
}

type HouseDealYT struct {
	*BaseHouseDeal
}

func NewHouseDealYT() *HouseDealYT {
	return &HouseDealYT{
		&BaseHouseDeal{
			"烟台", "yt", nil,
		},
	}
}

type HouseDealZS struct {
	*BaseHouseDeal
}

func NewHouseDealZS() *HouseDealZS {
	return &HouseDealZS{
		&BaseHouseDeal{
			"中山", "zs", nil,
		},
	}
}

type HouseDealZH struct {
	*BaseHouseDeal
}

func NewHouseDealZH() *HouseDealZH {
	return &HouseDealZH{
		&BaseHouseDeal{
			"珠海", "zh", nil,
		},
	}
}

type HouseDealSY struct {
	*BaseHouseDeal
}

func NewHouseDealSY() *HouseDealSY {
	return &HouseDealSY{
		&BaseHouseDeal{
			"沈阳", "sy", nil,
		},
	}
}

type HouseDealS struct {
	*BaseHouseDeal
}

func NewHouseDealS() *HouseDealS {
	return &HouseDealS{
		&BaseHouseDeal{
			"苏州", "s", nil,
		},
	}
}

type HouseDealLF struct {
	*BaseHouseDeal
}

func NewHouseDealLF() *HouseDealLF {
	return &HouseDealLF{
		&BaseHouseDeal{
			"廊坊", "lf", nil,
		},
	}
}

type HouseDealTY struct {
	*BaseHouseDeal
}

func NewHouseDealTY() *HouseDealTY {
	return &HouseDealTY{
		&BaseHouseDeal{
			"太原", "ty", nil,
		},
	}
}

type HouseDealHUI struct {
	*BaseHouseDeal
}

func NewHouseDealHUI() *HouseDealHUI {
	return &HouseDealHUI{
		&BaseHouseDeal{
			"惠州", "hui", nil,
		},
	}
}
