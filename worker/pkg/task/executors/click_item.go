package executors

import (
	"context"
	"encoding/json"
	"errors"
	"merge-api/shared/entity/task"
	"merge-api/worker/internal/service"
)

var UnableToClickNonclickableItem = errors.New("unable to click non-clickable item")

type ClickItemTaskExecutor struct {
	service           *service.CollectionBoard
	collectionService *service.Collection
}

func (n *ClickItemTaskExecutor) CanExecuteThisTask(t *task.Task) (bool, error) {
	return t.Type == task.ClickItem, nil
}

func (n *ClickItemTaskExecutor) Execute(t *task.Task) (any, error) {
	args := new(task.ClickItemTaskArgs)
	err := json.Unmarshal(t.Args, args)
	if err != nil {
		return nil, err
	}

	atCoords1, err1 := (*n.service).GetBoardByCoordinates(args.BoardID, args.W1, args.H1)
	if err1 != nil {
		return nil, err1
	}

	ctx := context.Background()
	isClickable, err := (*n.collectionService).IsItemClickable(ctx, atCoords1)
	if err != nil {
		return nil, err
	}
	if !isClickable {
		return nil, UnableToClickNonclickableItem
	}

	produceItem, err := (*n.collectionService).GetItemProduceResult(ctx, atCoords1)
	if err != nil {
		return nil, err
	}
	emptyCellW, emptyCellH, err := (*n.service).FindEmptyCell(args.BoardID)
	if err != nil {
		return nil, err
	}
	err = (*n.service).UpdateCell(args.BoardID, emptyCellW, emptyCellH, produceItem)
	if err != nil {
		return nil, err
	}
	//TODO: return a result
	return nil, nil
}

func NewClickItemTaskExecutor(boardService service.CollectionBoard, collection service.Collection) *ClickItemTaskExecutor {
	return &ClickItemTaskExecutor{
		service:           &boardService,
		collectionService: &collection,
	}
}
