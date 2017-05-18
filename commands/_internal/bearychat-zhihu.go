//usr/local/bin/go run $0 $@ ; exit
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	simplejson "github.com/bitly/go-simplejson"
)

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.zhihu.com/r/search?q=bearychat&type=content", nil)
	resp, err := client.Do(req)
	exitWithErrIfNeed(err)

	defer resp.Body.Close()
	json, err := simplejson.NewFromReader(resp.Body)
	exitWithErrIfNeed(err)

	htmls := json.GetPath("htmls")
	results := make(map[int]string)
	keys := []int{}

	for i := 0; i < len(htmls.MustArray([]interface{}{})); i++ {
		h := htmls.GetIndex(i).MustString("")
		if h == "" {
			os.Exit(1)
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(h))
		exitWithErrIfNeed(err)
		title := doc.Find(".title").Text()
		link := doc.Find("link").AttrOr("href", "")
		if link != "" {
			link = "https://www.zhihu.com" + link
		}

		comps := strings.Split(link, "/")
		answerIDStr := comps[len(comps)-1]
		answerID, err := strconv.Atoi(answerIDStr)
		exitWithErrIfNeed(err)

		author := doc.Find("span.author-link-line").Text()

		keys = append(keys, answerID)
		results[answerID] = fmt.Sprintf("%s\n%s: %s", title, author, link)
	}

	sort.Ints(keys)
	var buff bytes.Buffer
	for i := len(keys) - 1; i >= 0; i-- {
		buff.WriteString(results[keys[i]])
		if i != 0 {
			buff.WriteString("\n")
		}
	}
	fmt.Print(buff.String())
}

func exitWithErrIfNeed(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
