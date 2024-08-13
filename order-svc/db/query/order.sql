-- name: CreateOrder :one
INSERT INTO orders (order_ref_id, customer_id, username, product_name, status) 
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateOrder :one
UPDATE orders
SET 
    status = $1,
    total_amount = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE 
    order_ref_id = $3
RETURNING *;

-- name: CountByID :one
SELECT COUNT(*)
FROM orders
WHERE order_ref_id = $1;

-- name: FindOrderByID :one
SELECT * FROM orders WHERE id = $1 LIMIT 1;
