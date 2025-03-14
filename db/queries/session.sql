-- name: CreateSession :one
INSERT INTO session (
  id,
  username,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expires_at,
  created_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: SelectSession :one
SELECT * FROM session
WHERE username = $1;