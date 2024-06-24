package executors

import (
	"encoding/json"
	"github.com/google/uuid"
	"merge-api/shared/config"
	"merge-api/shared/entity/task"
	"merge-api/shared/pkg/database"
	collectionRepo "merge-api/worker/internal/repo/repos/collection"
	"merge-api/worker/internal/service"
	"merge-api/worker/internal/service/collection"
	"reflect"
	"testing"
	"time"
)

func prepareCollectionService() (error, service.Collection) {
	conf, err := config.NewConfigWithDiscover(nil)
	if err != nil {
		panic(err)
	}
	connection, err := database.NewDatabaseTestConnection(conf)
	if err != nil {
		panic(err)
	}
	repo := collectionRepo.NewCollectionRepo((*database.Database)(&connection))
	collectionService := collection.NewCollectionService(repo)
	return nil, collectionService
}

func TestMergeItemsTaskExecutor_Execute(t *testing.T) {
	err, connection, collectionBoardService := prepareDatabaseWithService()
	err, collectionService := prepareCollectionService()
	if err != nil {
		panic(err)
	}

	connection.SetUp("boards")
	connection.SetUp("board_cells")
	boardId := createBoard(collectionBoardService)

	type fields struct {
		service           *service.CollectionBoard
		collectionService *service.Collection
	}
	defaultFields := fields{
		service:           &collectionBoardService,
		collectionService: &collectionService,
	}
	type args struct {
		t *task.Task
	}

	bothEmptyArgs, err := json.Marshal(task.NewMergeItemsTaskArgs(boardId, 7, 7, 6, 6))
	if err != nil {
		panic(err)
	}

	sameCoordsArgs, err := json.Marshal(task.NewMergeItemsTaskArgs(boardId, 7, 7, 7, 7))
	if err != nil {
		panic(err)
	}

	firstEmptySecondItemArgs, err := json.Marshal(task.NewMergeItemsTaskArgs(boardId, 1, 2, 7, 7))
	if err != nil {
		panic(err)
	}

	firstItemSecondEmptyArgs, err := json.Marshal(task.NewMergeItemsTaskArgs(boardId, 7, 7, 1, 2))
	if err != nil {
		panic(err)
	}

	sameCollectionDifferentItemsArgs, err := json.Marshal(task.NewMergeItemsTaskArgs(boardId, 1, 2, 4, 2))
	if err != nil {
		panic(err)
	}

	differentCollectionsSameItemsArgs, err := json.Marshal(task.NewMergeItemsTaskArgs(boardId, 1, 2, 8, 9))
	if err != nil {
		panic(err)
	}

	validItemsArgs, err := json.Marshal(task.NewMergeItemsTaskArgs(boardId, 1, 2, 2, 4))
	if err != nil {
		panic(err)
	}
	validItemsCheckFunction := func() bool {
		at, err := collectionBoardService.GetBoardByCoordinates(boardId, 1, 2)
		if err != nil {
			panic(err)
		}
		cId, cIId := at.GetCollectionInfo()
		return cId == 1 && cIId == 2
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
			name:   "both cell are empty",
			fields: defaultFields,
			args: args{
				t: &task.Task{
					ID:                   1,
					UUID:                 uuid.New(),
					Type:                 task.MergeItems,
					Status:               task.Scheduled,
					Args:                 bothEmptyArgs,
					TimeCreated:          time.Time{},
					TimeStartedExecuting: time.Time{},
					TimeDoneExecuting:    time.Time{},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "same coordinates",
			fields: defaultFields,
			args: args{
				t: &task.Task{
					ID:                   1,
					UUID:                 uuid.New(),
					Type:                 task.MergeItems,
					Status:               task.Scheduled,
					Args:                 sameCoordsArgs,
					TimeCreated:          time.Time{},
					TimeStartedExecuting: time.Time{},
					TimeDoneExecuting:    time.Time{},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "first empty, second item",
			fields: defaultFields,
			args: args{
				t: &task.Task{
					ID:                   1,
					UUID:                 uuid.New(),
					Type:                 task.MergeItems,
					Status:               task.Scheduled,
					Args:                 firstEmptySecondItemArgs,
					TimeCreated:          time.Time{},
					TimeStartedExecuting: time.Time{},
					TimeDoneExecuting:    time.Time{},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "second empty, first item",
			fields: defaultFields,
			args: args{
				t: &task.Task{
					ID:                   1,
					UUID:                 uuid.New(),
					Type:                 task.MergeItems,
					Status:               task.Scheduled,
					Args:                 firstItemSecondEmptyArgs,
					TimeCreated:          time.Time{},
					TimeStartedExecuting: time.Time{},
					TimeDoneExecuting:    time.Time{},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "both item from the same collection but different",
			fields: defaultFields,
			args: args{
				t: &task.Task{
					ID:                   1,
					UUID:                 uuid.New(),
					Type:                 task.MergeItems,
					Status:               task.Scheduled,
					Args:                 sameCollectionDifferentItemsArgs,
					TimeCreated:          time.Time{},
					TimeStartedExecuting: time.Time{},
					TimeDoneExecuting:    time.Time{},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "both item from the different collections but items are equal",
			fields: defaultFields,
			args: args{
				t: &task.Task{
					ID:                   1,
					UUID:                 uuid.New(),
					Type:                 task.MergeItems,
					Status:               task.Scheduled,
					Args:                 differentCollectionsSameItemsArgs,
					TimeCreated:          time.Time{},
					TimeStartedExecuting: time.Time{},
					TimeDoneExecuting:    time.Time{},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "valid items to merge",
			fields: defaultFields,
			args: args{
				t: &task.Task{
					ID:                   1,
					UUID:                 uuid.New(),
					Type:                 task.MergeItems,
					Status:               task.Scheduled,
					Args:                 validItemsArgs,
					TimeCreated:          time.Time{},
					TimeStartedExecuting: time.Time{},
					TimeDoneExecuting:    time.Time{},
				},
			},
			checkFunction: &validItemsCheckFunction,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &MergeItemsTaskExecutor{
				service:           tt.fields.service,
				collectionService: tt.fields.collectionService,
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
