-- name: GetAllAccounts :many
SELECT * FROM accounts;
--

-- name: GetAllDevices :many
SELECT * FROM devices;
--

-- name: AddPlexDevice :one
INSERT INTO devices (id, name, last_seen)
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

-- name: DeleteDevice :exec
DELETE FROM devices
WHERE id = ?;
--

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = ?;
--
