package database

import (
	"database/sql"
	"github.com/WesJD/proxy-scraper/app/utils"
	"github.com/WesJD/proxy-scraper/app/config"
	"github.com/influxdata/influxdb/client/v2"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Sql *sql.DB
	Influx client.Client
	submitStatement *sql.Stmt
)

func Connect(config *config.Configuration) {
	influx, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Influx.Address,
		Username: config.Influx.Username,
		Password: config.Influx.Password,
	})
	Influx = influx
	utils.CheckError(err)

	sqlDb, err := sql.Open("mysql", config.DatabaseUrl)
	Sql = sqlDb
	utils.CheckError(err)

	//defaults
	Sql.Exec("CREATE TABLE IF NOT EXISTS proxies (ip_port CHAR(40) NOT NULL, checking BOOL NOT NULL, working BOOL NOT NULL, last_checked TIMESTAMP NOT NULL, UNIQUE (ip_port))")
	Sql.Exec(`
		DELIMITER //
		DROP PROCEDURE IF EXISTS matchProxies //

		CREATE PROCEDURE
  			matchProxies( amount INT, age TIMESTAMP )
		BEGIN
    		DECLARE _proxy CHAR(40);
    		DECLARE done INT;

    		DECLARE cur_proxies CURSOR FOR 
        		SELECT ip_port FROM proxies WHERE checking = 0 AND last_checked < age AND working = 1 LIMIT amount;
    		DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;

    		SELECT ip_port FROM proxies WHERE checking = 0 AND last_checked < age AND working = 1 LIMIT amount;
    
    		OPEN cur_proxies;

    		Reading_proxies: LOOP
        		FETCH NEXT FROM cur_proxies INTO _proxy;
        		IF done THEN
            		LEAVE Reading_proxies;
        		END IF;

        		UPDATE proxies SET checking = TRUE WHERE ip_port = _proxy;
    		END LOOP;
		END
		//

		DELIMITER ;
	`)

	stmt, err := Sql.Prepare("INSERT INTO proxies (ip_port, working) VALUES (?,?) ON DUPLICATE KEY UPDATE working=VALUES(working)")
	utils.CheckError(err)
	submitStatement = stmt
}