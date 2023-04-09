package requests

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Response struct {
	Status int
	Body   []byte
}

func getUserAgent() string {
	payload := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:102.0) Gecko/20100101 Firefox/102.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 12.4; rv:102.0) Gecko/20100101 Firefox/102.0",
		"Mozilla/5.0 (X11; Linux i686; rv:102.0) Gecko/20100101 Firefox/102.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 12_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/102.0 Mobile/15E148 Safari/605.1.15",
		"Mozilla/5.0 (iPod touch; CPU iPhone OS 12_4 like Mac OS X) AppleWebKit/604.5.6 (KHTML, like Gecko) FxiOS/102.0 Mobile/15E148 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 12.4; rv:102.0) Gecko/20100101 Firefox/102.0",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:102.0) Gecko/20100101 Firefox/102.0",
	}
	rand.Seed(time.Now().UnixNano())
	return payload[rand.Intn(len(payload))]
}

func makeRequest(urlx, method string, params io.Reader, headers map[string]string) (Response, error) {
	Client := &http.Client{Timeout: time.Second * 20}
	request, err := http.NewRequest(method, urlx, params)
	if err != nil {
		return Response{}, err
	}
	request.Header.Set("User-Agent", getUserAgent())
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	response, err := Client.Do(request)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return Response{}, err
	} else if err != nil {
		return Response{}, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Response{}, err
	}
	return Response{
		Status: response.StatusCode,
		Body:   body,
	}, nil
}

func Post(urlx string, params map[string]string, headers map[string]string) (Response, error) {
	p := url.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	return makeRequest(urlx, http.MethodPost, bytes.NewBufferString(p.Encode()), headers)
}

func Get(urlx string, params map[string]string, headers map[string]string) (Response, error) {
	parsed, _ := url.Parse(urlx)
	p := parsed.Query()
	for k, v := range params {
		p.Add(k, v)
	}
	parsed.RawQuery = p.Encode()
	urlx = parsed.String()
	return makeRequest(urlx, http.MethodGet, nil, headers)
}

// func PostJson(urlx string, params interface{}, headers map[string]string) (Response, error) {
// 	json, _ := json.Marshal(params)
// 	buf := bytes.NewBuffer(json)
// 	if headers == nil {
// 		headers = make(map[string]string)
// 	}
// 	headers["Content-Type"] = "application/json"
// 	return req(urlx, http.MethodPost, buf, headers)
// }
