package sites

import (
	"github.com/ddliu/go-httpclient"
	"fmt"
	"time"
	"encoding/json"
	"github.com/headzoo/surf/errors"
	"github.com/WesJD/proxy-scraper/utils"
)

type GetProxyListResponse struct {
	Ip   string
	Port int
}

type GetProxyList struct{}

func (s *GetProxyList) Check(url string) (result map[string]bool, err error) {
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

	if response.Port == 0 {
		err = errors.New("exceeded daily usage for get proxy list")
		return
	}
	address := fmt.Sprintf("%s:%d", response.Ip, response.Port)
	result[address] = utils.CheckProxy(address, url)

	return
}

func (s GetProxyList) WaitTime() time.Duration {
	return 5 * time.Second
}
