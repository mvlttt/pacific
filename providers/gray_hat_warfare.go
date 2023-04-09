package providers

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/pacific/requests"
)

type GrayHatWarfare struct {
	*Config
}

type grayHatWarfareFiles struct {
	Urls []struct {
		SubDomain string `json:"subDomain"`
		URL       string `json:"url"`
		Ext       string `json:"ext"`
		MimeType  string `json:"mimeType"`
		ShortURL  string `json:"shortUrl"`
	} `json:"urls"`
}

func NewGrayHatWarfare(config *Config) Provider {
	return &GrayHatWarfare{Config: config}
}

func (o *GrayHatWarfare) Fetch(domain string, results chan<- string) {
	var apiurl string = "https://shorteners.grayhatwarfare.com/api/v1/files/0/1000"
	apikey, err := o.getApiKey(o.Config.ApiKeys.Grayhatwarfare)
	if err != nil {
		return
	}
	params := map[string]string{
		"access_token": apikey,
		"keywords":     domain,
	}

	res, _ := requests.Get(apiurl, params, nil)
	var filesJson grayHatWarfareFiles
	json.Unmarshal(res.Body, &filesJson)

	for _, u := range filesJson.Urls {
		parsedUrl, _ := url.Parse(u.URL)
		if strings.Contains(parsedUrl.Host, domain) {
			results <- u.URL
		}
	}
}

func (o *GrayHatWarfare) getApiKey(apiKeys []string) (string, error) {
	if len(apiKeys) < 1 {
		return "", errors.New("Api keys not found")
	}
	//:TODO: add api limit control

	return apiKeys[0], nil
}
