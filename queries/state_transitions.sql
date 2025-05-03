-- name: RecordTransition :exec
INSERT INTO state_transitions (task_id, from_state, to_state, output)
VALUES (?, ?, ?, ?);

-- name: GetTaskState :one
SELECT to_state FROM state_transitions
WHERE task_id = ?
  AND to_state != '__error__'
ORDER BY created_at DESC
LIMIT 1;

-- name: GetHistory :many
SELECT * FROM state_transitions
WHERE task_id = ?
ORDER BY created_at ASC;

-- name: GetLastValidTransition :one
SELECT * FROM state_transitions
WHERE task_id = ?
  AND to_state != '__error__'
ORDER BY created_at DESC
LIMIT 1;

