-- name: GetAllBulbs :many
SELECT * FROM bulbs;
--

-- name: UpdateBulb :exec
UPDATE bulbs
SET ip = ?, updated_at = ?
WHERE mac = ?;

-- name: AddBulb :one
INSERT INTO bulbs (mac, created_at, updated_at, ip, name)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?
)
RETURNING *;
