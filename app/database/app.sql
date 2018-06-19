-- name: insert-proxies
INSERT INTO proxies (ip_port, working) VALUES (?, ?)
ON DUPLICATE KEY UPDATE working = VALUES(working);

-- name: update-proxy
UPDATE proxies SET working = ?, checking = FALSE WHERE ip_port = ?;

-- name: get-amount-working
SELECT COUNT(*) FROM proxies WHERE working = TRUE;