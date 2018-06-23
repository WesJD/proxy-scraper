package checking

import (
	"github.com/WesJD/proxy-scraper/app/config"
	"fmt"
	"github.com/WesJD/proxy-scraper/app/utils"
	"github.com/WesJD/proxy-scraper/app/database"
	"sync/atomic"
	"time"
)

func Start(config *config.Configuration, trueResponse string) {
	updateStatement, err := database.Sql.Prepare(database.Client, "update-proxy")
	utils.CheckError(err)

	for i := 0; i < config.Checking.Services; i++ {
		go func() {
			for {
				//cannot prepare a CALL statement... has to just stay here
				query := fmt.Sprintf("CALL matchProxies(%d, NOW())", config.Checking.PerRound)
				rows, err := database.Client.Query(query)
				utils.CheckError(err)
				for rows.Next() {
					var ipPort string
					err = rows.Scan(&ipPort)
					utils.CheckError(err)

					_, err := updateStatement.Exec(utils.CheckProxy(config.Static, trueResponse, ipPort), ipPort)
					utils.CheckError(err)

					atomic.AddInt64(&database.AmountChecked, 1)
				}
				rows.Close()
				time.Sleep(1)
			}
		}()
	}
}
