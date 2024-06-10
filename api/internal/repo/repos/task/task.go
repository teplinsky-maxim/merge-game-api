package task

import (
	"context"
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	task2 "merge-api/shared/entity/task"
	"merge-api/shared/pkg/database"
)

type TaskRepo struct {
	database *database.Database
}

func (r *TaskRepo) CreateTaskNewBoard(ctx context.Context, width, height uint) (task2.Task, error) {
	taskArgs := task2.NewNewBoardTaskArgs(width, height)
	marshalledArgs, err := taskArgs.MarshalJSON()
	if err != nil {
		return task2.Task{}, nil
	}
	return addTask(ctx, r, task2.MoveItem, marshalledArgs)
}

func (r *TaskRepo) CreateTaskMoveItem(ctx context.Context, boardId, w1, h1, w2, h2 uint) (task2.Task, error) {
	args := task2.NewMoveItemTaskArgs(boardId, w1, h1, w2, h2)
	marshalledArgs, err := args.MarshalJSON()
	if err != nil {
		return task2.Task{}, err
	}
	return addTask(ctx, r, task2.MoveItem, marshalledArgs)
}

func addTask(ctx context.Context, r *TaskRepo, taskType task2.Type, marshalledArgs json.RawMessage) (task2.Task, error) {
	tx, err := r.database.DB.Begin(ctx)
	if err != nil {
		return task2.Task{}, nil
	}

	taskUUID := uuid.New()
	stmt := sq.
		Insert("tasks").Columns("uuid", "type", "status", "args").
		Values(taskUUID, taskType, task2.Scheduled, marshalledArgs).
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
			result.Type = taskType
			result.Status = task2.Scheduled
			result.Args = marshalledArgs
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
