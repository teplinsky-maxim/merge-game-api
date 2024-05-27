package task

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	task2 "merge-api/internal/entity/task"
	"merge-api/pkg/board"
	"merge-api/pkg/database"
	"merge-api/pkg/task"
)

type TaskRepo struct {
	database *database.Database
}

func (r *TaskRepo) CreateTaskNewBoard(ctx context.Context, width, height board.SizeType) (task.IDType, error) {
	tx, err := r.database.DB.Begin(ctx)
	if err != nil {
		return [16]byte{}, err
	}

	taskUUID := uuid.New()
	taskArgs := task2.NewNewBoardTaskArgs(uint(width), uint(height))

	stmt := sq.
		Insert("tasks").Columns("uuid", "type", "status", "args").
		Values(taskUUID, task2.NewBoard, task2.Scheduled, taskArgs).
		Suffix("RETURNING id, uuid").
		PlaceholderFormat(sq.Dollar)

	query, args, err := stmt.ToSql()
	if err != nil {
		return [16]byte{}, err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return [16]byte{}, err
	}

	var result task2.Task
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return [16]byte{}, err
		} else {
			err = rows.Scan(&result.ID, &result.UUID)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return [16]byte{}, err
	}
	return task.IDType(result.UUID), err
}

func NewTaskRepo(database *database.Database) *TaskRepo {
	return &TaskRepo{
		database: database,
	}
}
