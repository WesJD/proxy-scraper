package database

import (
	"github.com/WesJD/proxy-scraper/app/utils"
	"database/sql"
)

var submitStatement *sql.Stmt

func makeStatements() {
	stmt, err := AppSql.Prepare(Sql, "insert-proxies")
	utils.CheckError(err)
	submitStatement = stmt
}


func SubmitProxies(proxies map[string]bool) {
	for proxy, working := range proxies {
		_, err := submitStatement.Exec(proxy, working)
		utils.CheckError(err)
	}
}

func ReportStats() {

}