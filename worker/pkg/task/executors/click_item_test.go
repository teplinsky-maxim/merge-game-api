package executors

import (
	"encoding/json"
	"github.com/google/uuid"
	"merge-api/shared/entity/task"
	"merge-api/worker/internal/service"
	"reflect"
	"testing"
	"time"
)

func TestClickItemTaskExecutor_Execute(t *testing.T) {
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

	emptyCellClickArgs, err := json.Marshal(task.NewClickItemTaskArgs(boardId, 5, 5))
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
			name:   "empty cell click",
			fields: defaultFields,
			args: args{
				t: &task.Task{
					ID:                   1,
					UUID:                 uuid.New(),
					Type:                 task.MergeItems,
					Status:               task.Scheduled,
					Args:                 emptyCellClickArgs,
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
			n := &ClickItemTaskExecutor{
				service:           tt.fields.service,
				collectionService: tt.fields.collectionService,
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
