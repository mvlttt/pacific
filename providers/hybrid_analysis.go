package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/pacific/requests"
)

type hybridAnalysisJobs struct {
	Result []struct {
		JobID string `json:"job_id"`
	} `json:"result"`
}

type hybridAnalysisHash []struct {
	SubmitName string `json:"submit_name"`
}

type HybridAnalysis struct {
	*Config
}

func NewHybridAnalysis(config *Config) Provider {
	return &HybridAnalysis{Config: config}
}

func (o *HybridAnalysis) Fetch(domain string, results chan<- string) {
	apikey, err := o.getApiKey(o.Config.ApiKeys.Grayhatwarfare)
	if err != nil {
		return
	}
	headers := map[string]string{"api-key": apikey, "User-Agent": "VxApi CLI Connector", "Content-Type": "application/x-www-form-urlencoded"}
	params := map[string]string{"domain": domain}
	const terms string = "https://www.hybrid-analysis.com/api/v2/search/terms"
	const summary string = "https://www.hybrid-analysis.com/api/v2/report/summary"

	res, _ := requests.Post(terms, params, headers)
	var jobs hybridAnalysisJobs
	json.Unmarshal(res.Body, &jobs)
	var wg sync.WaitGroup
	for k, job := range jobs.Result {
		wg.Add(1)
		go func(k int, job string) {
			defer wg.Done()
			params := map[string]string{
				fmt.Sprintf("hashes[%s]", strconv.Itoa(k)): job,
			}
			res, _ := requests.Post(summary, params, headers)
			hashs := hybridAnalysisHash{}
			json.Unmarshal(res.Body, &hashs)
			for _, hash := range hashs {
				results <- hash.SubmitName
			}
		}(k, job.JobID)

	}
	wg.Wait()
}

func (o *HybridAnalysis) getApiKey(apiKeys []string) (string, error) {
	if len(apiKeys) < 1 {
		return "", errors.New("Api keys not found")
	}
	//:TODO: add api limit control

	return apiKeys[0], nil
}
