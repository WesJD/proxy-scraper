package checkers

import (
		"time"
	"github.com/ddliu/go-httpclient"
	"encoding/json"
	"github.com/WesJD/proxy-scraper/app/utils"
	)

type PubProxyResponse struct {
	Data []PubProxyResponseData
}

type PubProxyResponseData struct {
	IpPort string
}

type PubProxy struct {}

func (s PubProxy) Check(url string, trueResponse string) (result *CheckResult, err error) {
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
	if err = json.Unmarshal([]byte(value), &response); err != nil {
		return
	}

	result = &CheckResult{}
	for _, proxy := range response.Data {
		if utils.CheckProxy(url, trueResponse, proxy.IpPort) {
			result.Passing++
			result.WorkingProxies = append(result.WorkingProxies, proxy.IpPort)
		} else {
			result.Failing++
		}
	}

	return
}

func (s PubProxy) WaitTime() time.Duration {
	return 1000 * 5 * time.Millisecond
}
