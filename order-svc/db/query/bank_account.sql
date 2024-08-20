-- name: CreateBankAccountRegistration :one
INSERT INTO bank_account_registration (customer_id, username, email, status, deposit)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateBankAccountRegistration :one
UPDATE bank_account_registration
SET
    status = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE
    username = $2 AND email = $3 AND customer_id = $4 RETURNING *;

-- name: FindBankAccountRegistrationByUsernameOrEmail :one
SELECT * FROM bank_account_registration
WHERE username = $1 OR email = $2 LIMIT 1;