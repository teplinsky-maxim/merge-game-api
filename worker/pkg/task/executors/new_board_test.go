package executors

import (
	"encoding/json"
	"github.com/google/uuid"
	"merge-api/shared/config"
	"merge-api/shared/entity/task"
	"merge-api/shared/pkg/database"
	"merge-api/worker/internal/repo/repos/collection/board"
	"merge-api/worker/internal/repo/repos/collection/redis_board"
	"merge-api/worker/internal/service"
	boardService "merge-api/worker/internal/service/board"
	"merge-api/worker/pkg/redis"
	"reflect"
	"testing"
	"time"
)

func prepareDatabaseWithService() (error, database.DatabaseTest, service.CollectionBoard) {
	conf, err := config.NewConfigWithDiscover(nil)
	if err != nil {
		panic(err)
	}
	connection, err := database.NewDatabaseTestConnection(conf)
	if err != nil {
		panic(err)
	}
	repo := board.NewBoardRepo((*database.Database)(&connection))
	redisClient, err := redis.NewTestRedis(conf)
	redisClient.Clear()
	if err != nil {
		panic(err)
	}
	redisRepo := redis_board.NewRedisBoardRepo((*redis.Redis)(&redisClient))
	collectionBoardService := boardService.NewCollectionBoardService(repo, redisRepo)
	return err, connection, collectionBoardService
}

func TestNewBoardTaskExecutor_Execute(t *testing.T) {
	err, connection, collectionBoardService := prepareDatabaseWithService()
	if err != nil {
		panic(err)
	}

	connection.SetUp("boards")

	type fields struct {
		service *service.CollectionBoard
	}
	defaultFields := fields{
		service: &collectionBoardService,
	}
	type args struct {
		t *task.Task
	}
	zeroWidthArgs, err := json.Marshal(task.NewNewBoardTaskArgs(0, 10))
	if err != nil {
		panic(err)
	}
	zeroHeightArgs, err := json.Marshal(task.NewNewBoardTaskArgs(10, 0))
	if err != nil {
		panic(err)
	}
	okArgs, err := json.Marshal(task.NewNewBoardTaskArgs(10, 9))
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		{
			name:   "zero width",
			fields: defaultFields,
			args: args{
				t: &task.Task{
					ID:                   1,
					UUID:                 uuid.New(),
					Type:                 task.NewBoard,
					Status:               task.Scheduled,
					Args:                 zeroWidthArgs,
					TimeCreated:          time.Time{},
					TimeStartedExecuting: time.Time{},
					TimeDoneExecuting:    time.Time{},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "zero height",
			fields: defaultFields,
			args: args{
				t: &task.Task{
					ID:                   1,
					UUID:                 uuid.New(),
					Type:                 task.NewBoard,
					Status:               task.Scheduled,
					Args:                 zeroHeightArgs,
					TimeCreated:          time.Time{},
					TimeStartedExecuting: time.Time{},
					TimeDoneExecuting:    time.Time{},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "ok",
			fields: defaultFields,
			args: args{
				t: &task.Task{
					ID:                   1,
					UUID:                 uuid.New(),
					Type:                 task.NewBoard,
					Status:               task.Scheduled,
					Args:                 okArgs,
					TimeCreated:          time.Time{},
					TimeStartedExecuting: time.Time{},
					TimeDoneExecuting:    time.Time{},
				},
			},
			want:    &task.NewBoardTaskResult{BoardID: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NewBoardTaskExecutor{
				service: tt.fields.service,
			}
			got, err := n.Execute(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Execute() got = %v, want %v", got, tt.want)
			}
		})
	}
}
