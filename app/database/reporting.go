package database

import "github.com/WesJD/proxy-scraper/app/utils"

func SubmitProxies(proxies map[string]bool) {
	for proxy, working := range proxies {
		_, err := submitStatement.Exec(proxy, working)
		utils.CheckError(err)
	}
}

func ReportStats() {

}