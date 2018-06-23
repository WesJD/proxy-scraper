package database

const (
	AppSql = `
	-- name: insert-proxies
	INSERT INTO proxies (ip_port, working) VALUES (?, ?)
	ON DUPLICATE KEY UPDATE working = VALUES(working);

	-- name: update-proxy
	UPDATE proxies SET working = ?, checking = FALSE WHERE ip_port = ?;

	-- name: get-amount-working
	SELECT COUNT(*) FROM proxies WHERE working = TRUE;

	-- name: set-all-not-checking
	UPDATE proxies SET checking = FALSE;
	`

	SetupSql = `
	-- name: setup-proxies
	CREATE TABLE IF NOT EXISTS proxies (
  	ip_port      CHAR(40)  NOT NULL,
  	checking     BOOL      NOT NULL,
  	working      BOOL      NOT NULL,
  	last_checked TIMESTAMP NOT NULL,
  	UNIQUE (ip_port)
	);

	-- name: create-procedure
	DELIMITER //
	DROP PROCEDURE IF EXISTS matchProxies //

	CREATE PROCEDURE
  		matchProxies( amount INT, age TIMESTAMP )
	BEGIN
    	START TRANSACTION;
    	SELECT ip_port FROM proxies WHERE checking = FALSE AND last_checked < age ORDER BY last_checked ASC LIMIT amount FOR UPDATE;
    	UPDATE proxies SET checking = TRUE WHERE checking = FALSE AND last_checked < age ORDER BY last_checked ASC LIMIT amount;
    	COMMIT;
	END
	//

	DELIMITER ;
	`
)
