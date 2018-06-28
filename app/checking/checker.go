package checking

import (
	"github.com/WesJD/proxy-scraper/app/config"
	"fmt"
	"github.com/WesJD/proxy-scraper/app/utils"
	"github.com/WesJD/proxy-scraper/app/database"
	"sync/atomic"
	"time"
)

func Start(trueResponse string) {
	updateSuccessStatement, err := database.Sql.Prepare(database.Client, "update-proxy-working")
	proxyFailedQuery := "CALL proxyFailed('%s')"
	utils.CheckError(err)

	for i := 0; i < config.Values.Checking.Services; i++ {
		go func() {
			for {
				//cannot prepare a CALL statement... has to just stay here
				query := fmt.Sprintf("CALL matchProxies(%d, NOW())", config.Values.Checking.PerRound)
				rows, err := database.Client.Query(query)
				utils.CheckError(err)
				for rows.Next() {
					var ipPort string
					err = rows.Scan(&ipPort)
					utils.CheckError(err)

					checkResult := utils.CheckProxy(trueResponse, ipPort)

					if checkResult {
						_, err := updateSuccessStatement.Exec(ipPort)
						utils.CheckError(err)
					} else {
						failRows, err := database.Client.Query(fmt.Sprintf(proxyFailedQuery, ipPort))
						utils.CheckError(err)

						var consecFails int64

						for failRows.Next() {
							err = failRows.Scan(&consecFails)
							utils.CheckError(err)
						}

						failRows.Close()
					}

					atomic.AddInt64(&database.AmountChecked, 1)
				}
				rows.Close()
				time.Sleep(1)
			}
		}()
	}
}
