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
  BEGIN;
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