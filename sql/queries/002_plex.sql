-- name: AddPlexPlayer :one
INSERT INTO players (uuid, title, public_address)
VALUES (
    ?,
    ?,
    ?
)
RETURNING *;
--

-- name: AddPlexAccount :one
INSERT INTO accounts (id, title)
VALUES (
    ?,
    ?
)
RETURNING *;
--

-- name: DeletePlayer :exec
DELETE FROM players
WHERE uuid = ?;
--

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = ?;
--
