package providers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pacific/requests"
)

type Alienvault struct {
	*Config
}

type alienvaultOtxDomain struct {
	URLList []struct {
		URL string `json:"url"`
	} `json:"url_list"`
	HasNext bool `json:"has_next"`
}

func NewAlienvault(config *Config) Provider {
	return &Alienvault{Config: config}
}

func (o *Alienvault) Fetch(domain string, results chan<- string) {
	for _, p := range []string{"domain", "hostname"} {
		var url string = fmt.Sprintf("https://otx.alienvault.com/api/v1/indicator/%s/%s/url_list/?limit=50", p, domain)

		var page int = 1
		for {
			params := map[string]string{
				"page": strconv.Itoa(page),
			}
			res, _ := requests.Get(url, params, nil)
			var otxDomain alienvaultOtxDomain
			json.Unmarshal(res.Body, &otxDomain)
			for _, result := range otxDomain.URLList {
				results <- result.URL
			}

			if !otxDomain.HasNext {
				break
			}
			page += 1
		}
	}
}
