-- these procedures are to be created manually, placed here for convenience

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