package utils

import (
	"github.com/ddliu/go-httpclient"
	"github.com/headzoo/surf/errors"
)

func CheckProxy(proxyIp string, url string) bool {
	res, err := httpclient.
		Begin().
		WithOption(httpclient.OPT_PROXY, proxyIp).
		Get(url)
	if err != nil {
		return false
	}

	value, err := res.ToString()
	if err != nil {
		return false
	}

	return value == GetCached(url)
}

func CheckProxyAndReason(proxyIp string, url string) (success bool, err error) {
	res, err := httpclient.
		Begin().
		WithOption(httpclient.OPT_PROXY, proxyIp).
		Get(url)
	if err != nil {
		return false, err
	}

	value, err := res.ToString()
	if err != nil {
		return false, err
	}
	err = errors.New("Didn't match page with trueResponse")
	return value == GetCached(url), err
}