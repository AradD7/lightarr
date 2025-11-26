-- name: AddRule :one
INSERT INTO rules (id, condition, action)
VALUES (
    ?,
    ?,
    ?
)
RETURNING *;
--

-- name: DeleteRule :exec
DELETE FROM rules
WHERE id = ?;
--

-- name: GetAllRules :many
SELECT * FROM rules;
--
