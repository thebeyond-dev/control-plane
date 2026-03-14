-- name: UpdateSubscription :exec
UPDATE "user"
SET devices = $1,
    bandwidth = $2,
    subscription_expires_at = subscription_expires_at + ($3 * INTERVAL '1 day')
WHERE id = $4;

-- name: DeactivateSubscription :exec
UPDATE "user" SET
    devices = 0,
    bandwidth = 0,
    subscription_expires_at = NULL
WHERE id = $1;