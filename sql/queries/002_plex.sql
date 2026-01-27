-- name: GetAllAccounts :many
SELECT * FROM accounts;
--

-- name: GetAllDevices :many
SELECT * FROM devices;
--

-- name: AddPlexDevice :one
INSERT INTO devices (id, name, product)
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

-- name: GetPlexDeviceById :one
SELECT * FROM devices
WHERE id = ?;
--

-- name: GetPlexAccountById :one
SELECT * FROM accounts
WHERE id = ?;
--
