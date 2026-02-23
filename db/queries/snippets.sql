-- name: ListSnippets :many
SELECT id, title, expires FROM snippets ORDER BY title;

-- name: CreateSnippet :one
INSERT INTO snippets (title, content, created, expires)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteSnippet :exec
DELETE FROM snippets
WHERE id = $1;

-- name: GetSnippetNotExpired :one
SELECT id, title, content, created, expires FROM snippets
WHERE expires > NOW() AND id = $1;

-- name: GetActiveSnippetsLimit10 :many
SELECT id, title, content, created, expires FROM snippets
    WHERE expires > NOW() ORDER BY id DESC LIMIT 10;