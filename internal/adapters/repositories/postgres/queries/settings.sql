-- name: SetUserLanguageCode :exec
UPDATE "user"
SET language_code = $2
WHERE id = $1;

-- name: SetUserCurrencyCode :exec
UPDATE "user"
SET currency_code = $2
WHERE id = $1;