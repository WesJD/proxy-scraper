package utils

import (
	"github.com/ddliu/go-httpclient"
	"github.com/WesJD/proxy-scraper/app/config"
)

func CheckProxy(trueResponse string, proxyIp string) bool {
	res, err := httpclient.
		Begin().
		WithOption(httpclient.OPT_PROXY, proxyIp).
		WithOption(httpclient.OPT_TIMEOUT_MS, config.Values.Scraping.TimeoutMs).
		Get(config.Values.Scraping.Static)
	if err != nil {
		return false
	}

	value, err := res.ToString()
	if err != nil {
		return false
	}

	return value == trueResponse
}