-- name: ListProducts :many
SELECT * FROM products;

-- name: FindProductById :one
SELECT * FROM products WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders (
  customer_id
) VALUES ($1) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price_cents)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateProductQuantity :one
UPDATE products
SET quantity = $2
WHERE id = $1
RETURNING *;

-- name: CreateProduct :one
INSERT INTO products (name, price_in_cents, quantity) 
VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (name, username, password)
VALUES ($1, $2, $3) RETURNING *;

-- name: FindUserById :one
SELECT * FROM users WHERE id = $1;

-- name: FindUserByUsername :one
SELECT * FROM users WHERE username = $1;