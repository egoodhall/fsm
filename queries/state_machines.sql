-- name: CreateStateMachine :one
INSERT INTO state_machines (name)
VALUES (?)
ON CONFLICT (name) DO UPDATE SET id = id
RETURNING id;

-- name: GetStateMachine :one
SELECT * FROM state_machines
WHERE id = ?;
