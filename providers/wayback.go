package providers

import (
	"fmt"
	"strings"

	"github.com/pacific/requests"
)

type Wayback struct {
	*Config
}

func NewWayback(config *Config) Provider {
	return &Wayback{Config: config}
}

func (o *Wayback) Fetch(domain string, results chan<- string) {
	var url string = fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=%s/&matchType=domain&collapse=urlkey&showResumeKey=false&fl=original", domain)
	res, _ := requests.Get(url, nil, nil)
	for _, url := range strings.Split(string(res.Body), "\n") {
		results <- url
	}
}
