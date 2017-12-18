package beary_commands

import (
	"github.com/bearyinnovative/lili/commands/house"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	"github.com/bearyinnovative/lili/notifier/bearychat"
	. "github.com/bearyinnovative/lili/util"
)

func GetAllDealCommands() (results []CommandType) {
	data := [][]string{
		[]string{"北京", "bj"},
		[]string{"上海", "sh"},
		[]string{"广州", "gz"},
		[]string{"深圳", "sz"},
		[]string{"天津", "tj"},
		[]string{"成都", "cd"},
		[]string{"南京", "nj"},
		[]string{"杭州", "hz"},
		[]string{"青岛", "qd"},
		[]string{"大连", "dl"},
		[]string{"厦门", "xm"},
		[]string{"武汉", "wh"},
		[]string{"重庆", "cq"},
		[]string{"长沙", "cs"},
		[]string{"西安", "xa"},
		[]string{"济南", "jn"},
		[]string{"石家庄", "sjz"},
		[]string{"东莞", "dg"},
		[]string{"佛山", "fs"},
		[]string{"合肥", "hf"},
		[]string{"烟台", "yt"},
		[]string{"中山", "zs"},
		[]string{"珠海", "zh"},
		[]string{"沈阳", "sy"},
		[]string{"苏州", "s"},
		[]string{"廊坊", "lf"},
		[]string{"太原", "ty"},
		[]string{"惠州", "hui"},
	}

	for _, d := range data {
		var notifiers []NotifierType
		if d[1] == "sz" {
			notifiers = szNotifiers()
		}

		results = append(results, &house.BaseHouseDeal{
			d[0], d[1], notifiers,
		})
	}

	return
}

func szNotifiers() []NotifierType {
	var n NotifierType
	n, err := bearychat.NewRTMNotifier("4f2dda2fa66a0d1fc575d341cca4eda6", "=bwG5y")
	if LogIfErr(err) {
		n = BCChannelNotifier("house_info")
	}
	return []NotifierType{
		n,
	}
}
