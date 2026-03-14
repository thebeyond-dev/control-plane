-- name: GetDeviceByPublicKey :one
SELECT pubkey, node_id, name FROM device WHERE user_id = $1 AND pubkey = $2;

-- name: ListDevices :many
SELECT pubkey, node_id, name FROM device WHERE user_id = $1;

-- name: CreateDevice :exec
INSERT INTO device(user_id, pubkey, node_id, name)
VALUES ($1, $2, $3, $4);

-- name: DeleteDevice :exec
DELETE FROM device WHERE user_id = $1 AND pubkey = $2;