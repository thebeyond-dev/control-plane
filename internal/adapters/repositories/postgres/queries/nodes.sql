-- name: GetNode :one
SELECT id, address FROM node WHERE id = $1;

-- name: ListNodes :many
SELECT id, address FROM node;