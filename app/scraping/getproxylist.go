package scraping

import (
	"github.com/ddliu/go-httpclient"
	"fmt"
	"time"
	"encoding/json"
	"github.com/WesJD/proxy-scraper/app/utils"
)

type GetProxyListResponse struct {
	Ip   string
	Port int
}

type GetProxyList struct{}

func (s GetProxyList) Check(url string, trueResponse string) (result map[string]bool, err error) {
	res, err := httpclient.
		Begin().
		Get("https://api.getproxylist.com/proxy?protocol[]=http&anonymity[]=high%20anonymity&anonymity[]=anonymous")
	if err != nil {
		return
	}

	var response GetProxyListResponse
	value, err := res.ToString()
	if err != nil {
		return
	}

	if err = json.Unmarshal([]byte(value), &response); err != nil {
		return
	}

	result = make(map[string]bool)

	address := fmt.Sprintf("%s%d", response.Ip, response.Port)
	fmt.Println("I got an address", address)
	result[address] = utils.CheckProxy(url, trueResponse, address)

	return
}

func (s GetProxyList) WaitTime() time.Duration {
	return 1000 * 5 * time.Millisecond
}
