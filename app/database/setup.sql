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
    SELECT ip_port FROM proxies WHERE checking = FALSE AND last_checked < age LIMIT amount FOR UPDATE;
    UPDATE proxies SET checking = TRUE WHERE checking = FALSE AND last_checked < age LIMIT amount;
    COMMIT;
END
//

DELIMITER ;