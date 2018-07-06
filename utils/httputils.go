package utils

import (
	"github.com/ddliu/go-httpclient"
)

var (
	pageCache = make(map[string]string)
)

func CachePage(url string) string {
	response := GetPage(url)
	pageCache[url] = response
	return response
}

func GetCached(url string) (value string) {
	value = pageCache[url]
	if value == "" { //nothing
		return CachePage(url)
	} else {
		return value
	}
}

func GetPage(url string) (trueResponse string) {
	res, err := httpclient.Get(url)
	CheckError(err)
	trueResponse, err = res.ToString()
	CheckError(err)
	return
}
