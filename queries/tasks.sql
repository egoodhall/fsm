-- name: CreateTask :one
INSERT INTO tasks (fsm_id, event)
VALUES (?, ?)
RETURNING *;

-- name: ListTasks :many
SELECT * FROM tasks
ORDER BY id ASC;
