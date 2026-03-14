-- name: GetUserByIdentity :one
SELECT u.id,
    u.devices,
    u.bandwidth,
    u.subscription_expires_at,
    u.language_code,
    u.currency_code,
    u.referrer_id,
    u.referral_balance,
    u.referral_commission_rate,
    u.referrals_count,
    u.created_at
FROM "user" u
JOIN identity i ON i.user_id = u.id
WHERE i.provider = $1 AND i.provider_id = $2
LIMIT 1;

-- name: CreateUser :exec
INSERT INTO "user"(
    id, devices, bandwidth, subscription_expires_at,
    language_code, currency_code, referrer_id
) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: IncrementReferralsCount :exec
UPDATE "user"
SET referrals_count = referrals_count + 1
WHERE id = $1;

-- name: AddBonusDays :exec
UPDATE "user"
SET subscription_expires_at = CASE
    WHEN subscription_expires_at > CURRENT_TIMESTAMP
    THEN subscription_expires_at + INTERVAL '5 days'
    ELSE CURRENT_TIMESTAMP + INTERVAL '5 days'
END
WHERE id = $1;

-- name: CreateIdentity :exec
INSERT INTO identity(provider, provider_id, user_id)
VALUES ($1, $2, $3);