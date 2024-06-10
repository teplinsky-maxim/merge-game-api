package task

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	task2 "merge-api/shared/entity/task"
	"merge-api/shared/pkg/database"
)

type TaskRepo struct {
	database *database.Database
}

func (r *TaskRepo) CreateTaskNewBoard(ctx context.Context, width, height uint) (task2.Task, error) {
	tx, err := r.database.DB.Begin(ctx)
	if err != nil {
		return task2.Task{}, nil
	}

	taskUUID := uuid.New()
	taskArgs := task2.NewNewBoardTaskArgs(width, height)
	taskArgsMarshalled, err := taskArgs.MarshalJSON()
	if err != nil {
		return task2.Task{}, nil
	}

	stmt := sq.
		Insert("tasks").Columns("uuid", "type", "status", "args").
		Values(taskUUID, task2.NewBoard, task2.Scheduled, taskArgsMarshalled).
		Suffix("RETURNING id, time_created").
		PlaceholderFormat(sq.Dollar)

	query, args, err := stmt.ToSql()
	if err != nil {
		return task2.Task{}, nil
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return task2.Task{}, nil
	}

	var result task2.Task
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return task2.Task{}, nil
		} else {
			err = rows.Scan(
				&result.ID,
				&result.TimeCreated,
			)
			if err != nil {
				return task2.Task{}, err
			}
			result.UUID = taskUUID
			result.Type = task2.NewBoard
			result.Status = task2.Scheduled
			result.Args = taskArgsMarshalled
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return task2.Task{}, nil
	}
	return result, err
}

func NewTaskRepo(database *database.Database) *TaskRepo {
	return &TaskRepo{
		database: database,
	}
}
