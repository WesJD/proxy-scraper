package utils

import (
	"github.com/ddliu/go-httpclient"
	"fmt"
)

func CheckProxy(url string, trueResponse string, proxyIp string) bool {
	res, err := httpclient.
		Begin().
		WithOption(httpclient.OPT_PROXY, proxyIp).
		Get(url)
	if err != nil {
		fmt.Println(proxyIp, "nar", err.Error())
		return false
	}

	value, err := res.ToString()
	if err != nil {
		return false
	}

	return value == trueResponse
}