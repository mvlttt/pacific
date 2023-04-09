package providers

import (
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/pacific/requests"
)

type commonCrawlApi []struct {
	CdxAPI string `json:"cdx-api"`
}

type commonCrawlArchive struct {
	URL string `json:"url"`
}
type CommonCrawl struct {
	*Config
}

func NewCommonCrawl(config *Config) Provider {
	return &CommonCrawl{Config: config}
}

func (o *CommonCrawl) Fetch(domain string, results chan<- string) {
	const url string = "https://index.commoncrawl.org/collinfo.json"
	params := map[string]string{
		"output":    "json",
		"matchType": "domain",
		"url":       domain,
	}
	res, err := requests.Get(url, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	var apiData commonCrawlApi
	json.Unmarshal(res.Body, &apiData)
	var wg sync.WaitGroup
	for _, api := range apiData {
		wg.Add(1)
		var badStatus bool = false
		go func(api string) {
			defer wg.Done()
			res, err := requests.Get(api, params, nil)
			if err != nil {
				for i := 0; i < 3; i++ {
					res, err = requests.Get(api, params, nil)
					if err == nil {
						badStatus = false
						break
					} else {
						badStatus = true
					}
				}
			}
			if badStatus {
				close(make(chan bool))
			}
			var archive commonCrawlArchive
			for _, jsonData := range strings.Split(string(res.Body), "\n") {
				json.Unmarshal([]byte(jsonData), &archive)
				results <- archive.URL
			}

		}(api.CdxAPI)
	}
	wg.Wait()

}
