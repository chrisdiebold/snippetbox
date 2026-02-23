-- name: ListSnippets :many
SELECT id, title, expires FROM snippets ORDER BY title;

-- name: CreateSnippet :one
INSERT INTO snippets (title, content, created, expires)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteSnippet :exec
DELETE FROM snippets
WHERE id = $1;