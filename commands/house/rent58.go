package house

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"
)

type Rent58 struct {
	Province  string
	District  string
	RoomNum   int
	Query     string
	Notifiers []NotifierType
}

func (c *Rent58) GetName() string {
	return fmt.Sprintf("58-rent-%s-%s-%d-%s", c.Province, c.District, c.RoomNum, c.Query)
}

func (c *Rent58) GetInterval() time.Duration {
	return time.Minute*45 + time.Minute*time.Duration(rand.Intn(15))
}

func (z *Rent58) Fetch() (results []*Item, err error) {
	room := ""
	if z.RoomNum > 0 {
		room = fmt.Sprintf("/j%d", z.RoomNum)
	}
	district := ""
	if z.District != "" {
		district = "/" + z.District
	}
	path := fmt.Sprintf("http://%s.58.com%s/zufang/0%s/?key=",
		z.Province, district, room)
	log.Printf("[%s] fetching: %s%s\n", z.GetName(), path, z.Query)
	path += url.PathEscape(z.Query)

	doc, err := respDoc(path)
	if LogIfErr(err) {
		return
	}

	doc.Find("body > div.mainbox > div.main > div.content > div.listBox > ul > li").Each(func(i int, s *goquery.Selection) {
		item := z.createItem(s)
		if item != nil {
			results = append(results, item)
		}
	})

	return
}

/*
   <li logr="p_0_33210274865675_32537449201458_0_0_sortid:568894957@postdate:1515550957000@ses:onlyruletitleranker^0" sortid="1515550957000" class="">
               <div class="img_list">
                   <a href="http://gz.58.com/zufang/32537449201458x.shtml?from=1-list-0" tongji_label="listclick" onclick="clickLog('from=fcpc_zflist_gzcount');" target="_blank">
                       <img lazy_src="http://pic7.58cdn.com.cn/dwater/fang/small/n_v2e168fa311ab74da6800a64a260eaeb19.jpg?wt=%40%E9%BB%84%E5%85%88%E7%94%9F&amp;ws=2d11969a6d8ec1c688f9a007f3ac9d4a&amp;w=294&amp;h=220&amp;crop=1" src="http://pic7.58cdn.com.cn/dwater/fang/small/n_v2e168fa311ab74da6800a64a260eaeb19.jpg?wt=%40%E9%BB%84%E5%85%88%E7%94%9F&amp;ws=2d11969a6d8ec1c688f9a007f3ac9d4a&amp;w=294&amp;h=220&amp;crop=1">
                   </a>
                   <span class="picNum">12 图</span>
               </div>
               <div class="des">
                   <h2>
                       <a href="http://gz.58.com/zufang/32537449201458x.shtml?from=1-list-0" tongji_label="listclick" onclick="clickLog('from=fcpc_zflist_gzcount');" target="_blank">
                           <b>保利心语</b><b>花园</b> 2室2厅1卫                    </a>
                                       </h2>
                   <p class="room">2室2厅1卫                    &nbsp;&nbsp;&nbsp;&nbsp;80㎡</p>
                   <p class="add">
                       <a href="/zhujiangxincheng/zufang/" onclick="clickLog('from=fcpc_list_gz_biaoti_shangquan')">珠江新城</a>
                       &nbsp;&nbsp;
                                               <a href="http://gz.58.com/xiaoqu/baolixinyuhuayuan/chuzu/" target="_blank" onclick="clickLog('from=fcpc_list_gz_biaoti_xiaoqu')"><b>保利心语</b><b>花园</b>租房</a>
                                                                   <em></em>距离5号线猎德地铁站444米                                    </p>
                                       <p class="geren">
                           <span>来自个人房源</span>：黄先生                    </p>
                               </div>
               <div class="listliright">
                   <div class="sendTime">
                       11小时前                </div>
                   <div class="money">
                       <b>7800</b>元/月                </div>
               </div>
               <div class="listline"></div>
           </li>
*/
func (z *Rent58) createItem(s *goquery.Selection) *Item {
	// filter 置顶房源
	if len(s.Find("div.des h2 a.dingico_a").Nodes) > 0 {
		return nil
	}

	reg := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	areaReg := regexp.MustCompile(`\d+㎡`)

	money := s.Find("div.money").Text()
	money = strings.TrimSpace(money)

	ref := s.Find("div.des h2 a").AttrOr("href", "")
	if ref == "" {
		return nil
	}

	timeStr := s.Find("div.sendTime").Text()
	timeStr = strings.TrimSpace(timeStr)
	timeStr = reg.ReplaceAllString(timeStr, " ")

	// "2017-10-04", ignore parse time error
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	created, _ := time.ParseInLocation("2006-01-02", timeStr, loc)

	des := s.Find("div.des h2 a").Text()
	des = strings.TrimSpace(des)
	des = reg.ReplaceAllString(des, " ")

	// 有的时候会返回一些不太相关的... 原因不明
	if !strings.Contains(des, z.Query) {
		log.Printf("[%s] `%s` doesn't contains `%s`\n", z.GetName(), des, z.Query)
		return nil
	}

	area := s.Find("div.des p.room").Text()
	area = areaReg.FindString(area)
	if area == "" {
		return nil
	}

	position := s.Find("div.des p.add").Text()
	position = strings.TrimSpace(position)
	position = reg.ReplaceAllString(position, " ")

	// [保利心语花园 2室1厅1卫 81㎡](http://gz.58.com/zufang/32273403139404x.shtml) 6000元/月 2017-11-04
	// 珠江新城 保利心语花园租房
	itemDesc := fmt.Sprintf("[%s %s](%s) %s %s\n%s",
		des, area, ref, money, timeStr,
		position)
	// fmt.Println(itemDesc)

	return &Item{
		Name:       z.GetName(),
		Identifier: z.GetName() + "-" + ref,
		Desc:       itemDesc,
		Ref:        ref,
		Key:        money,
		Created:    created,
		Notifiers:  z.Notifiers,
	}
}

func respDoc(path string) (*goquery.Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}
	// Display Results
	// fmt.Println("response Status : ", resp.Status)
	// fmt.Println("response Headers : ", resp.Header)

	doc, err := goquery.NewDocumentFromReader(reader)
	return doc, err
}
