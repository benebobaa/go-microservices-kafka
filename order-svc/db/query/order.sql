-- name: CreateOrder :one
INSERT INTO orders (ref_id, customer_id, username, product_id, quantity, status)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateOrder :one
UPDATE orders
SET 
    status = $1,
    amount = $2,
    quantity = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE 
    ref_id = $4
RETURNING *;

-- name: CountByID :one
SELECT COUNT(*)
FROM orders
WHERE ref_id = $1;

-- name: FindOrderByID :one
SELECT * FROM orders WHERE id = $1 LIMIT 1;

-- name: FindOrderByRefID :one
SELECT * FROM orders WHERE ref_id = $1 LIMIT 1;

-- name: FindOrdersByUsername :many
SELECT * FROM orders
WHERE username = $1
ORDER BY created_at  DESC;
