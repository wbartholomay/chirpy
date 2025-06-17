-- name: GetRefreshTokenByUser :one
SELECT * FROM refresh_tokens
WHERE user_id = $1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    NOW()+ INTERVAL '60 days',
    NULL
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT * FROM users
WHERE users.id = (
    SELECT refresh_tokens.user_id FROM refresh_tokens
    WHERE token = $1
);