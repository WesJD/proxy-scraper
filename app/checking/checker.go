package checking

import (
	"github.com/WesJD/proxy-scraper/app/config"
	"time"
	"fmt"
	"github.com/WesJD/proxy-scraper/app/utils"
	"github.com/WesJD/proxy-scraper/app/database"
)

// have failure amount
func Start(config *config.Configuration, trueResponse string) {
	updateStatement, err := database.AppSql.Prepare(database.Sql, "update-proxy")
	utils.CheckError(err)

	for i := 0; i < config.Checking.Services; i++ {
		go func() {
			for {
				//cannot prepare a CALL statement... has to just stay here
				query := fmt.Sprintf("CALL matchProxies(%d, NOW() - INTERVAL %s)", config.Checking.PerRound, config.Checking.OlderThan)
				rows, err := database.Sql.Query(query)
				utils.CheckError(err)
				for rows.Next() {
					var ipPort string
					err = rows.Scan(&ipPort)
					utils.CheckError(err)

					_, err := updateStatement.Exec(utils.CheckProxy(config.Static, trueResponse, ipPort), ipPort)
					utils.CheckError(err)
				}
				rows.Close()
				time.Sleep(config.Checking.EveryMs * time.Millisecond)
			}
		}()
	}
}
