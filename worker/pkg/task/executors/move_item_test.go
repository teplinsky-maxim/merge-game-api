package executors

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"merge-api/shared/entity/task"
	"merge-api/worker/internal/repo/repos/collection/redis_board"
	"merge-api/worker/internal/service"
	"merge-api/worker/pkg"
	"reflect"
	"testing"
	"time"
)

func createBoard(service service.CollectionBoard) uint {
	_, boardId, err := service.CreateBoard(10, 5)
	if err != nil {
		panic(err)
	}

	item1 := pkg.NewCollectionItemImpl(1, 1)
	item2 := pkg.NewCollectionItemImpl(1, 2)
	err = service.UpdateCell(boardId, 1, 2, &item1)
	if err != nil {
		panic(err)
	}
	err = service.UpdateCell(boardId, 2, 4, &item1)
	if err != nil {
		panic(err)
	}
	err = service.UpdateCell(boardId, 4, 2, &item2)
	if err != nil {
		panic(err)
	}

	return boardId
}

func TestMoveItemTaskExecutor_Execute(t *testing.T) {
	err, connection, collectionBoardService := prepareDatabaseWithService()
	if err != nil {
		panic(err)
	}

	connection.SetUp("boards")
	connection.SetUp("board_cells")
	boardId := createBoard(collectionBoardService)

	type fields struct {
		service *service.CollectionBoard
	}
	defaultFields := fields{
		service: &collectionBoardService,
	}
	type args struct {
		t *task.Task
	}

	firstEmptySecondItemArgs, err := json.Marshal(task.NewMoveItemTaskArgs(boardId, 1, 1, 1, 2))
	if err != nil {
		panic(err)
	}
	firstEmptySecondItemCheckFunction := func() bool {
		//became item{1, 1}
		at1, err := collectionBoardService.GetBoardByCoordinates(boardId, 1, 1)
		if err != nil {
			panic(err)
		}
		//became empty
		_, err = collectionBoardService.GetBoardByCoordinates(boardId, 1, 2)
		if !errors.Is(err, redis_board.BoardCellEmptyError) {
			return false
		}
		cId, cIId := at1.GetCollectionInfo()
		return cId == 1 && cIId == 1
	}

	firstItemSecondEmptyArgs, err := json.Marshal(task.NewMoveItemTaskArgs(boardId, 2, 4, 2, 2))
	if err != nil {
		panic(err)
	}
	firstItemSecondEmptyCheckFunction := func() bool {
		_, err = collectionBoardService.GetBoardByCoordinates(boardId, 2, 4)
		if !errors.Is(err, redis_board.BoardCellEmptyError) {
			return false
		}
		at2, err := collectionBoardService.GetBoardByCoordinates(boardId, 2, 2)
		if err != nil {
			panic(err)
		}
		cId, cIId := at2.GetCollectionInfo()
		return cId == 1 && cIId == 1
	}

	bothEmpty, err := json.Marshal(task.NewMoveItemTaskArgs(boardId, 3, 5, 7, 8))
	if err != nil {
		panic(err)
	}
	bothEmptyCheckFunction := func() bool {
		_, err = collectionBoardService.GetBoardByCoordinates(boardId, 3, 5)
		if !errors.Is(err, redis_board.BoardCellEmptyError) {
			return false
		}
		_, err = collectionBoardService.GetBoardByCoordinates(boardId, 7, 8)
		return errors.Is(err, redis_board.BoardCellEmptyError)
	}

	sameCoords, err := json.Marshal(task.NewMoveItemTaskArgs(boardId, 1, 1, 1, 1))
	if err != nil {
		panic(err)
	}

	sameCoordsFreeCell, err := json.Marshal(task.NewMoveItemTaskArgs(boardId, 9, 9, 9, 9))
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name          string
		fields        fields
		args          args
		checkFunction *func() bool
		want          any
		wantErr       bool
	}{
		{
			name:   "first cell is empty, second one has item",
			fields: defaultFields,
			args: args{
				t: &task.Task{
					ID:                   1,
					UUID:                 uuid.New(),
					Type:                 task.MoveItem,
					Status:               task.Scheduled,
					Args:                 firstEmptySecondItemArgs,
					TimeCreated:          time.Time{},
					TimeStartedExecuting: time.Time{},
					TimeDoneExecuting:    time.Time{},
				},
			},
			checkFunction: &firstEmptySecondItemCheckFunction,
		},
		{
			name:   "second cell is empty, first one has item",
			fields: defaultFields,
			args: args{
				t: &task.Task{
					ID:                   1,
					UUID:                 uuid.New(),
					Type:                 task.MoveItem,
					Status:               task.Scheduled,
					Args:                 firstItemSecondEmptyArgs,
					TimeCreated:          time.Time{},
					TimeStartedExecuting: time.Time{},
					TimeDoneExecuting:    time.Time{},
				},
			},
			checkFunction: &firstItemSecondEmptyCheckFunction,
		},
		{
			name:   "both empty",
			fields: defaultFields,
			args: args{
				t: &task.Task{
					ID:                   1,
					UUID:                 uuid.New(),
					Type:                 task.MoveItem,
					Status:               task.Scheduled,
					Args:                 bothEmpty,
					TimeCreated:          time.Time{},
					TimeStartedExecuting: time.Time{},
					TimeDoneExecuting:    time.Time{},
				},
			},
			checkFunction: &bothEmptyCheckFunction,
		},
		{
			name:   "same coordinates item cell",
			fields: defaultFields,
			args: args{
				t: &task.Task{
					ID:                   1,
					UUID:                 uuid.New(),
					Type:                 task.MoveItem,
					Status:               task.Scheduled,
					Args:                 sameCoords,
					TimeCreated:          time.Time{},
					TimeStartedExecuting: time.Time{},
					TimeDoneExecuting:    time.Time{},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "same coordinates free cell",
			fields: defaultFields,
			args: args{
				t: &task.Task{
					ID:                   1,
					UUID:                 uuid.New(),
					Type:                 task.MoveItem,
					Status:               task.Scheduled,
					Args:                 sameCoordsFreeCell,
					TimeCreated:          time.Time{},
					TimeStartedExecuting: time.Time{},
					TimeDoneExecuting:    time.Time{},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &MoveItemTaskExecutor{
				service: tt.fields.service,
			}
			got, err := n.Execute(tt.args.t)
			if tt.checkFunction != nil {
				checkResult := (*tt.checkFunction)()
				if checkResult != true {
					t.Error("checkFunction() returned false")
				}
			} else {
				if (err != nil) != tt.wantErr {
					t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Execute() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
