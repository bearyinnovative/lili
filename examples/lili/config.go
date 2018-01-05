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

type ConfigNotifier struct {
	Type string `yaml:"type"`

	// bearychat.incoming
	URL       string `yaml:"url"`
	ToUser    string `yaml:"to_user,omitempty"`
	ToChannel string `yaml:"to_channel,omitempty"`

	// telegram/telegram.media
	Token string `yaml:"token"`
	// `@channel_name` or integer id as string: `-123456`
	ChatID    string `yaml:"chat_id"`
	ParseMode string `yaml:"parse_mode,omitempty"`
}

func (cn *ConfigNotifier) toNotifierType() NotifierType {
	switch cn.Type {
	case "bearychat.incoming":
		if cn.URL == "" {
			log.Fatal("can't find")
			return nil
		}
		return &IncomingNotifier{
			URL:       cn.URL,
			ToUser:    cn.ToUser,
			ToChannel: cn.ToChannel,
		}
	case "telegram":
		return &telegram.Notifier{
			Token:     cn.Token,
			ChatID:    cn.ChatID,
			ParseMode: cn.ParseMode,
		}
	case "telegram.media":
		return telegram.NewMediaNotifier(cn.Token, cn.ChatID)
	default:
		log.Fatal("type unknown:", cn.Type)
		return nil
	}
}

type Config struct {
	Zhihu []struct {
		Keywords  []string          `yaml:"keywords"`
		Notifiers []*ConfigNotifier `yaml:notifiers,omitempty`
	} `yaml:"zhihu"`

	V2EX []struct {
		Keywords  []string          `yaml:"keywords"`
		Notifiers []*ConfigNotifier `yaml:notifiers,omitempty`
	} `yaml:"v2ex"`

	Instagram []struct {
		Tags      []string          `yaml:"tags,omitempty"`
		Notifiers []*ConfigNotifier `yaml:notifiers,omitempty`
		Usernames []string          `yaml:"usernames,omitempty"`
	} `yaml:"instagram"`

	Hackernews []struct {
		Name            string            `yaml:"name"`
		Keywords        []string          `yaml:"keywords,omitempty"`
		Notifiers       []*ConfigNotifier `yaml:notifiers,omitempty`
		MinScore        int               `yaml:"min_score,omitempty"`
		MinCommentCount int               `yaml:"min_comment_count,omitempty"`
	} `yaml:"hackernews"`

	Douban []struct {
		ID        string            `yaml:"id"`
		Notifiers []*ConfigNotifier `yaml:notifiers,omitempty`
	} `yaml:"douban"`

	HouseDeal []struct {
		Name      string            `yaml:"name"`
		Notifiers []*ConfigNotifier `yaml:notifiers,omitempty`
	} `yaml:"house_deal"`

	House []struct {
		Name        string `yaml:"name"`
		Subscribers []struct {
			MinPrice  int               `yaml:"min_price"`
			Notifiers []*ConfigNotifier `yaml:notifiers,omitempty`
		} `yaml:"subscribers,omitempty"`
	} `yaml:"house"`

	LocalBitcoin []struct {
		Currency  string            `yaml:"currency"`
		Interval  int               `yaml:"interval"`
		Notifiers []*ConfigNotifier `yaml:notifiers,omitempty`
	} `yaml:"localbitcoin"`

	CoinMarket []struct {
		Currency  string            `yaml:"currency"`
		Interval  int               `yaml:"interval"`
		Notifiers []*ConfigNotifier `yaml:notifiers,omitempty`
	} `yaml:"coinmarket"`

	Reddit []struct {
		Subreddits  []string          `yaml:"subreddits"`
		Interval    int               `yaml:"interval"`
		MinUpsRatio float64           `yaml:"min_ups_ratio"`
		ImageOnly   bool              `yaml:"image_only"`
		MinScore    int               `yaml:"min_score"`
		Notifiers   []*ConfigNotifier `yaml:notifiers,omitempty`
	} `yaml:"reddit"`
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
			Notifiers: toNotifierTypes(c.Notifiers),
		})
	}

	for _, c := range config.Zhihu {
		for _, keyword := range c.Keywords {
			if keyword == "" {
				continue
			}

			results = append(results, &BaseZhihu{
				Notifiers: toNotifierTypes(c.Notifiers),
				Query:     keyword,
			})
		}
	}

	var hnSubs = []*HackerNewsSubscriber{}
	for _, c := range config.Hackernews {
		minScore := c.MinScore
		minCommentCount := c.MinCommentCount
		keywords := c.Keywords

		hnSubs = append(hnSubs, &HackerNewsSubscriber{
			Name:      c.Name,
			Notifiers: toNotifierTypes(c.Notifiers),
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
	if len(hnSubs) > 0 {
		results = append(results, &HackerNews{Subscribers: hnSubs})
	}

	for _, c := range config.HouseDeal {
		notifiers := toNotifierTypes(c.Notifiers)
		deal, err := house.NewHouseDeal(c.Name, notifiers)
		if err != nil {
			log.Println("Error:", err)
			continue
		}

		results = append(results, deal)
	}

	for _, c := range config.House {
		subs := []*house.HouseSubscriber{}
		for _, s := range c.Subscribers {
			notifiers := toNotifierTypes(s.Notifiers)
			subs = append(subs, &house.HouseSubscriber{
				Notifiers: notifiers,
				ShouldNotify: func(hi *house.HouseItem) bool {
					return true
				},
			})
		}

		cmd, err := house.NewHouseSecondHand(c.Name, subs)
		if err != nil {
			log.Println("Error:", err)
			continue
		}

		// log.Printf("appending %s with %d subscribers", c.Name, len(subs))
		results = append(results, cmd)
	}

	for _, c := range config.Instagram {
		for _, tag := range c.Tags {
			if tag == "" {
				continue
			}

			notifiers := toNotifierTypes(c.Notifiers)
			results = append(results, NewTagInstagram(notifiers, tag))
		}

		for _, username := range c.Usernames {
			if username == "" {
				continue
			}

			notifiers := toNotifierTypes(c.Notifiers)
			results = append(results, NewUserInstagram(notifiers, username))
		}
	}

	for _, c := range config.V2EX {
		for _, keyword := range c.Keywords {
			if keyword == "" {
				continue
			}

			results = append(results, &BaseV2EX{
				Notifiers: toNotifierTypes(c.Notifiers),
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
			Notifiers: toNotifierTypes(c.Notifiers),
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
			Notifiers: toNotifierTypes(c.Notifiers),
		})
	}

	for _, c := range config.Reddit {
		if c.Interval <= 0 {
			c.Interval = 5
		}
		if c.MinUpsRatio < 0 || c.MinUpsRatio > 1 {
			c.MinUpsRatio = 0
		}

		for _, subreddit := range c.Subreddits {
			if subreddit == "" {
				log.Println("need subreddit")
				continue
			}

			results = append(results, &Reddit{
				Subreddit:   subreddit,
				Interval:    c.Interval,
				ImageOnly:   c.ImageOnly,
				MinUpsRatio: c.MinUpsRatio,
				MinScore:    c.MinScore,
				Notifiers:   toNotifierTypes(c.Notifiers),
			})
		}
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

func toNotifierTypes(notifiers []*ConfigNotifier) []NotifierType {
	results := []NotifierType{}

	for i := range notifiers {
		n := notifiers[i].toNotifierType()
		if n != nil {
			results = append(results, n)
		}
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
