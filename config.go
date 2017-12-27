package main

import (
	"io/ioutil"
	"log"
	"strings"

	. "github.com/bearyinnovative/lili/commands"
	"github.com/bearyinnovative/lili/commands/house"
	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/notifier/bearychat"
	"github.com/bearyinnovative/lili/notifier/telegram"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Zhihu []struct {
		Keywords          []string             `yaml:"keywords"`
		Notifiers         []*IncomingNotifier  `yaml:"notifiers,omitempty"`
		TelegramNotifiers []*telegram.Notifier `yaml:"telegram_notifiers,omitempty"`
	} `yaml:"zhihu"`

	V2EX []struct {
		Keywords          []string             `yaml:"keywords"`
		Notifiers         []*IncomingNotifier  `yaml:"notifiers,omitempty"`
		TelegramNotifiers []*telegram.Notifier `yaml:"telegram_notifiers,omitempty"`
	} `yaml:"v2ex"`

	Instagram []struct {
		Tags              []string             `yaml:"tags,omitempty"`
		Notifiers         []*IncomingNotifier  `yaml:"notifiers,omitempty"`
		TelegramNotifiers []*telegram.Notifier `yaml:"telegram_notifiers,omitempty"`
		Usernames         []string             `yaml:"usernames,omitempty"`
	} `yaml:"instagram"`

	Hackernews []struct {
		Name              string               `yaml:"name"`
		Keywords          []string             `yaml:"keywords,omitempty"`
		Notifiers         []*IncomingNotifier  `yaml:"notifiers,omitempty"`
		TelegramNotifiers []*telegram.Notifier `yaml:"telegram_notifiers,omitempty"`
		MinScore          int                  `yaml:"min_score,omitempty"`
		MinCommentCount   int                  `yaml:"min_comment_count,omitempty"`
	} `yaml:"hackernews"`

	Douban []struct {
		ID                string               `yaml:"id"`
		Notifiers         []*IncomingNotifier  `yaml:"notifiers,omitempty"`
		TelegramNotifiers []*telegram.Notifier `yaml:"telegram_notifiers,omitempty"`
	} `yaml:"douban"`

	HouseDeal []struct {
		Name              string               `yaml:"name"`
		ShortName         string               `yaml:"short_name"`
		Notifiers         []*IncomingNotifier  `yaml:"notifiers,omitempty"`
		TelegramNotifiers []*telegram.Notifier `yaml:"telegram_notifiers,omitempty"`
	} `yaml:"house_deal"`

	LocalBitcoin []struct {
		Currency          string               `yaml:"currency"`
		Interval          int                  `yaml:"interval"`
		Notifiers         []*IncomingNotifier  `yaml:"notifiers,omitempty"`
		TelegramNotifiers []*telegram.Notifier `yaml:"telegram_notifiers,omitempty"`
	} `yaml:"localbitcoin"`

	CoinMarket []struct {
		Currency          string               `yaml:"currency"`
		Interval          int                  `yaml:"interval"`
		Notifiers         []*IncomingNotifier  `yaml:"notifiers,omitempty"`
		TelegramNotifiers []*telegram.Notifier `yaml:"telegram_notifiers,omitempty"`
	} `yaml:"coinmarket"`
}

func (config *Config) ToCommandTypes() []CommandType {
	results := []CommandType{}

	// douban
	for _, c := range config.Douban {
		if c.ID == "" {
			log.Println("can't find douban id:", c)
			continue
		}

		results = append(results, &DoubanStatus{
			ID:        c.ID,
			Notifiers: toNotifierTypes(c.Notifiers, c.TelegramNotifiers),
		})
	}

	for _, c := range config.Zhihu {
		for _, keyword := range c.Keywords {
			if keyword == "" {
				continue
			}

			results = append(results, &BaseZhihu{
				Notifiers: toNotifierTypes(c.Notifiers, c.TelegramNotifiers),
				Query:     keyword,
			})
		}
	}

	for _, c := range config.Hackernews {
		if c.Name == "" {
			log.Println("can't find name for hackernews:", c)
			continue
		}

		minScore := c.MinScore
		minCommentCount := c.MinCommentCount
		keywords := c.Keywords

		results = append(results, &BaseHackerNews{
			Notifiers: toNotifierTypes(c.Notifiers, c.TelegramNotifiers),
			Name:      c.Name,
			ShouldNotify: func(item *HNItem) bool {
				if minScore > 0 && item.Score < minScore {
					return false
				}

				if minCommentCount > 0 && item.Comments < minCommentCount {
					return false
				}

				if len(keywords) > 0 && !checkContains(item.Title, keywords) {
					return false
				}

				return true
			},
		})
	}

	for _, c := range config.HouseDeal {
		if c.Name == "" || c.ShortName == "" {
			log.Println("can't find names for house deal:", c)
			continue
		}

		results = append(results, &house.BaseHouseDeal{
			CityName:      c.Name,
			CityShortName: c.ShortName,
			Notifiers:     toNotifierTypes(c.Notifiers, c.TelegramNotifiers),
		})
	}

	for _, c := range config.Instagram {
		for _, tag := range c.Tags {
			if tag == "" {
				continue
			}

			notifiers := toNotifierTypes(c.Notifiers, c.TelegramNotifiers)
			results = append(results, NewTagInstagram(notifiers, tag))
		}

		for _, username := range c.Usernames {
			if username == "" {
				continue
			}

			notifiers := toNotifierTypes(c.Notifiers, c.TelegramNotifiers)
			results = append(results, NewUserInstagram(notifiers, username))
		}
	}

	for _, c := range config.V2EX {
		for _, keyword := range c.Keywords {
			if keyword == "" {
				continue
			}

			results = append(results, &BaseV2EX{
				Notifiers: toNotifierTypes(c.Notifiers, c.TelegramNotifiers),
				Query:     keyword,
			})
		}
	}

	for _, c := range config.LocalBitcoin {
		if c.Currency == "" {
			log.Println("can't find currency for LocalBitcoin:", c)
			continue
		}

		// default interval 5 min
		if c.Interval <= 0 {
			c.Interval = 5
		}

		results = append(results, &BaseLBBuyOnline{
			Currency:  c.Currency,
			Interval:  c.Interval,
			Notifiers: toNotifierTypes(c.Notifiers, c.TelegramNotifiers),
		})
	}

	for _, c := range config.CoinMarket {
		if c.Currency == "" {
			log.Println("can't find currency for CoinMarket:", c)
			continue
		}

		// default interval 5 min
		if c.Interval <= 0 {
			c.Interval = 5
		}

		results = append(results, &CoinMarket{
			Currency:  c.Currency,
			Interval:  c.Interval,
			Notifiers: toNotifierTypes(c.Notifiers, c.TelegramNotifiers),
		})
	}

	return results
}

func checkContains(title string, keywords []string) bool {
	lowerTitle := strings.ToLower(title)

	for _, key := range keywords {
		if strings.Contains(lowerTitle, key) {
			return true
		}
	}

	return false
}

func toNotifierTypes(notifiers []*IncomingNotifier, telegramNotifiers []*telegram.Notifier) []NotifierType {
	len1, len2 := len(notifiers), len(telegramNotifiers)
	results := make([]NotifierType, len1+len2)

	for i := range notifiers {
		results[i] = notifiers[i]
	}

	for i := range telegramNotifiers {
		results[i+len1] = telegramNotifiers[i]
	}

	return results
}

func NewConfigFromFile(path string) (*Config, error) {
	var config *Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
