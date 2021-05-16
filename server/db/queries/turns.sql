-- name: ListTurns :many
SELECT id, person_id, team_id, date, created_at
FROM turns
WHERE team_id = ?
ORDER BY date DESC;

-- name: GetTurn :one
SELECT id, person_id, team_id, date, created_at
FROM turns
WHERE id = ?
AND team_id = ?
LIMIT 1;

-- name: CreateTurn :execresult
INSERT INTO turns (person_id, team_id, date)
VALUES ( ?, ?, ? );

-- name: UpdateTurn :execresult
UPDATE turns
SET person_id = ?, team_id = ?, date = ?
WHERE id = ?;

-- name: DeleteTurn :exec
DELETE FROM turns
WHERE id = ?
AND team_id = ?;