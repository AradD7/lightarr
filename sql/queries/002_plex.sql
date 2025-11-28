-- name: GetAllAccounts :many
SELECT * FROM accounts;
--

-- name: GetAllPlayers :many
SELECT * FROM players;
--

-- name: AddPlexPlayer :one
INSERT INTO players (id, name, last_seen)
VALUES (
    ?,
    ?,
    ?
)
RETURNING *;
--

-- name: AddPlexAccount :one
INSERT INTO accounts (id, title, thumb)
VALUES (
    ?,
    ?,
    ?
)
RETURNING *;
--

-- name: DeletePlayer :exec
DELETE FROM players
WHERE id = ?;
--

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = ?;
--
