package commands

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/dghubble/oauth1"
)

const apiBase = "https://api.flickr.com/services/rest/"

type FlickrClient struct {
	consumer *oauth1.Config
	user     *oauth1.Token
	client   *http.Client
}

func NewFlickrClient(consumerKey string, consumerSecret string, token string, tokenSecret string) *FlickrClient {
	c := &FlickrClient{
		consumer: oauth1.NewConfig(consumerKey, consumerSecret),
		user:     oauth1.NewToken(token, tokenSecret),
	}
	return c
}

func (c *FlickrClient) Get(method string) (*FlickrResp, error) {
	return c.GetWithParams(method, url.Values{})
}

func (c *FlickrClient) GetWithParams(method string, params url.Values) (*FlickrResp, error) {
	params.Set("method", method)
	return getResponse(c.GetHttpClient().Get(createRequestURI(apiBase, params)))
}

// Retrieve the underlying HTTP client
func (c *FlickrClient) GetHttpClient() *http.Client {
	if c.consumer == nil {
		panic("Consumer credentials are not set")
	}
	if c.client == nil {
		c.client = c.consumer.Client(context.TODO(), c.user)
		c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	return c.client
}

func createRequestURI(base string, params url.Values) string {
	setupDefaultParams(params)

	base += "?" + params.Encode()
	return base
}

func setupDefaultParams(params url.Values) {
	params.Set("format", "json")
	params.Set("nojsoncallback", "1")
	if params.Get("count") == "" {
		params.Set("count", "50")
	}
	if params.Get("extras") == "" {
		params.Set("extras", "url_l,media")
	}
}

func getResponse(resp *http.Response, err error) (*FlickrResp, error) {
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, errors.New(resp.Status)
	}

	var response *FlickrResp
	defer resp.Body.Close()

	// data, err := ioutil.ReadAll(resp.Body)

	// log.Println(string(data))
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type FlickrResp struct {
	Photos struct {
		Photo []struct {
			ID       string `json:"id"`
			Owner    string `json:"owner"`
			Username string `json:"username"`
			Title    string `json:"title"`
			Media    string `json:"media"`
			URLL     string `json:"url_l"`
			HeightL  string `json:"height_l"`
			WidthL   string `json:"width_l"`
		} `json:"photo"`
	} `json:"photos"`
	Stat string `json:"stat"`
}
