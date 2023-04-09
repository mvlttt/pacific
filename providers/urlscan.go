package providers

import (
	"encoding/json"

	"sync"

	"github.com/pacific/requests"
)

type urlscanResults struct {
	Results []struct {
		Page struct {
			URL string `json:"url"`
		} `json:"page"`
		Result string `json:"result"`
	} `json:"results"`
	Total int `json:"total"`
}

type urlscanData struct {
	Lists struct {
		Urls []string `json:"urls"`
	} `json:"lists"`
}

type Urlscan struct {
	*Config
}

func NewUrlscan(config *Config) Provider {
	return &Urlscan{Config: config}
}

func (o *Urlscan) Fetch(domain string, results chan<- string) {
	const url string = "https://urlscan.io/api/v1/search/"
	params := map[string]string{
		"q":    "domain:" + domain,
		"sort": "date:desc",
		"size": "1000",
	}
	res, _ := requests.Get(url, params, nil)
	var usresults urlscanResults
	json.Unmarshal(res.Body, &usresults)
	if usresults.Total == 0 {
		println("[-] Urlscan results <nil>")
	}
	var wg sync.WaitGroup
	for _, result := range usresults.Results {
		wg.Add(1)
		go func(u, r string) {
			defer wg.Done()
			results <- u
			res, _ := requests.Get(r, nil, nil)
			var data urlscanData
			json.Unmarshal(res.Body, &data)
			for _, url := range data.Lists.Urls {
				results <- url
			}
		}(result.Page.URL, result.Result)
	}
	wg.Wait()
}
