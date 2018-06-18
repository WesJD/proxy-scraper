package scraping

import (
	"time"
	"github.com/ddliu/go-httpclient"
	"encoding/json"
	"github.com/WesJD/proxy-scraper/app/utils"
	"fmt"
)

type PubProxyResponse struct {
	Data []PubProxyResponseData
}

type PubProxyResponseData struct {
	IpPort string
}

type PubProxy struct{}

func (s *PubProxy) Check(url string, trueResponse string) (result map[string]bool, err error) {
	res, err := httpclient.
		Begin().
		Get("http://pubproxy.com/api/proxy?limit=20&level=anonymous&level=elite")
	if err != nil {
		return
	}

	var response PubProxyResponse
	value, err := res.ToString()
	if err != nil {
		return
	}
	fmt.Println(value)
	if err = json.Unmarshal([]byte(value), &response); err != nil {
		return
	}

	result = make(map[string]bool)

	for _, proxy := range response.Data {
		address := proxy.IpPort
		result[address] = utils.CheckProxy(url, trueResponse, address)
	}

	return
}

func (s PubProxy) WaitTime() time.Duration {
	return 1000 * 5 * time.Millisecond
}
