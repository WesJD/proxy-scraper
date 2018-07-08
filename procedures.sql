-- these procedures are to be created manually, placed here for convenience

DELIMITER //
DROP PROCEDURE IF EXISTS matchProxies //

CREATE PROCEDURE
  	matchProxies( amount INT, age TIMESTAMP, max_consec_fails INT )
BEGIN
    START TRANSACTION;
   	    SELECT ip_port FROM proxies WHERE last_checked < age AND consec_fails <= max_consec_fails ORDER BY last_checked ASC LIMIT amount FOR UPDATE;
   	    UPDATE proxies SET last_checked = NOW() WHERE last_checked < age AND consec_fails <= max_consec_fails ORDER BY last_checked ASC LIMIT amount;
   	COMMIT;
END
//

DELIMITER ;
