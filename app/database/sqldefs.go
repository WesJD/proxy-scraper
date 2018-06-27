package database

const (
	AppSql = `
	-- name: insert-proxies
	INSERT INTO proxies (ip_port, working, consec_fails) VALUES (?, ?, 0)
	ON DUPLICATE KEY UPDATE working = VALUES(working), consec_fails = 0;

	-- name: update-proxy-working
	UPDATE proxies SET working = TRUE, checking = FALSE, consec_fails = 0 WHERE ip_port = ?;

	-- name: get-amount-working
	SELECT COUNT(*) FROM proxies WHERE working = TRUE;

	-- name: set-all-not-checking
	UPDATE proxies SET checking = FALSE;
	`

	SetupSql = `
	-- name: setup-proxies
	CREATE TABLE IF NOT EXISTS proxies (
  	ip_port      CHAR(40)         NOT NULL,
  	checking     BOOL             NOT NULL,
  	working      BOOL             NOT NULL,
  	last_checked TIMESTAMP        NOT NULL,
	consec_fails INTEGER UNSIGNED NOT NULL,
  	UNIQUE (ip_port)
	);

	-- name: create-match-proxies-procedure
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

	-- name: create-update-fail-procedure
	DELIMITER //
	DROP PROCEDURE IF EXISTS proxyFailed //

	CREATE PROCEDURE
  		proxyFailed( in_ip_port CHAR(40) )
	BEGIN
    	START TRANSACTION;
		UPDATE proxies SET working = FALSE, checking = FALSE, consec_fails = consec_fails + 1 WHERE ip_port = in_ip_port;
    	SELECT consec_fails FROM proxies WHERE ip_port = in_ip_port FOR UPDATE;
    	COMMIT;
	END
	//

	DELIMITER ;
	`
)
