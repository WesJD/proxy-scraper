package scraping

import (
	"time"
	"github.com/ddliu/go-httpclient"
	"encoding/json"
	"github.com/WesJD/proxy-scraper/app/utils"
	"strings"
	"github.com/headzoo/surf/errors"
)

type PubProxyResponse struct {
	Data []PubProxyResponseData
}

type PubProxyResponseData struct {
	IpPort string
}

type PubProxy struct{}

func (s *PubProxy) Check(trueResponse string) (result map[string]bool, err error) {
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
	if strings.Contains(value, "reached the maximum") {
		err = errors.New("reached the maximum amount of requests")
		return
	}
	if err = json.Unmarshal([]byte(value), &response); err != nil {
		return
	}

	result = make(map[string]bool)
	for _, proxy := range response.Data {
		address := proxy.IpPort
		result[address] = utils.CheckProxy(trueResponse, address)
	}

	return
}

func (s *PubProxy) WaitTime() time.Duration {
	return time.Minute
}
