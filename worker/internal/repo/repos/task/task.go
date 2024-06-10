package task

import (
	"context"
	"encoding/json"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"merge-api/shared/entity/task"
	"merge-api/shared/pkg/database"
	"time"
)

var (
	NotFound = errors.New("task not found")
)

type Repo struct {
	database *database.Database
}

func (r *Repo) UpdateTask(ctx context.Context, taskId task.IDType, status task.Status, result any) error {
	tx, err := r.database.DB.Begin(ctx)
	if err != nil {
		return err
	}

	stmt := sq.
		Update("tasks").
		Where(sq.Eq{"uuid": uuid.UUID(taskId).String()}).
		Set("status", status).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	if status == task.Running {
		stmt = stmt.Set("time_started_executing", time.Now())
	} else if status == task.Failed || status == task.Done {
		stmt = stmt.Set("time_done_executing", time.Now())
	}

	if status == task.Done {
		marshalledResult, err := json.Marshal(result)
		if err != nil {
			return err
		}
		stmt = stmt.Set("result", marshalledResult)
	}

	query, args, err := stmt.ToSql()
	if err != nil {
		return err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return err
	}

	id := 0
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			_ = tx.Rollback(ctx)
			return err
		} else {
			err = rows.Scan(&id)
		}
	}
	if id == 0 {
		_ = tx.Rollback(ctx)
		return NotFound
	}
	err = tx.Commit(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return nil
}

func NewTaskRepo(database *database.Database) *Repo {
	return &Repo{
		database: database,
	}
}
