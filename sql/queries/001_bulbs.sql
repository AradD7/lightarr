-- name: GetAllBulbs :many
SELECT * FROM bulbs;
--

-- name: UpdateBulbIp :exec
UPDATE bulbs
SET ip = ?, updated_at = ?
WHERE mac = ?;

-- name: AddBulb :one
INSERT INTO bulbs (mac, created_at, updated_at, ip, name, is_reachable)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
)
RETURNING *;
--

-- name: UpdateBulbName :exec
UPDATE bulbs
SET name = ?, updated_at = ?
WHERE mac = ?;
--

-- name: DeleteBulb :exec
DELETE FROM bulbs
where mac = ?;
--
