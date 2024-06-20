package executors

import (
	"context"
	"encoding/json"
	"errors"
	"merge-api/shared/entity/task"
	"merge-api/worker/internal/service"
)

var UnableToMergeDifferentCollections = errors.New("unable to merge different collections")
var UnableToMergeDifferentCollectionItems = errors.New("unable to merge different items")
var UnableToMergeNonmergeableItem = errors.New("unable to merge non-mergeable items")

type MergeItemsTaskExecutor struct {
	service           *service.CollectionBoard
	collectionService *service.Collection
}

func (n *MergeItemsTaskExecutor) CanExecuteThisTask(t *task.Task) (bool, error) {
	return t.Type == task.MergeItems, nil
}

func (n *MergeItemsTaskExecutor) Execute(t *task.Task) (any, error) {
	args := new(task.MergeItemsTaskArgs)
	err := json.Unmarshal(t.Args, args)
	if err != nil {
		return nil, err
	}

	if args.W1 == args.W2 && args.H1 == args.H2 {
		return nil, CoordinatesAreTheSameError
	}

	atCoords1, err1 := (*n.service).GetBoardByCoordinates(args.BoardID, args.W1, args.H1)
	if err1 != nil {
		return nil, err1
	}
	atCoords2, err2 := (*n.service).GetBoardByCoordinates(args.BoardID, args.W2, args.H2)
	if err2 != nil {
		return nil, err2
	}
	cId1, cIId1 := atCoords1.GetCollectionInfo()
	cId2, cIId2 := atCoords2.GetCollectionInfo()

	if cId1 != cId2 {
		return nil, UnableToMergeDifferentCollections
	}
	if cIId1 != cIId2 {
		return nil, UnableToMergeDifferentCollectionItems
	}
	ctx := context.Background()
	isMergeable, err := (*n.collectionService).IsItemMergeable(ctx, atCoords1)
	if err != nil {
		return nil, err
	}
	if !isMergeable {
		return nil, UnableToMergeNonmergeableItem
	}

	nextItem, err := (*n.collectionService).GetNextCollectionItem(ctx, atCoords1)
	if err != nil {
		return nil, err
	}
	err = (*n.service).UpdateCell(args.BoardID, args.W1, args.W2, nextItem)
	if err != nil {
		return nil, err
	}
	err = (*n.service).ClearCell(args.BoardID, args.W2, args.H2)
	if err != nil {
		return nil, err
	}
	//TODO: return a result
	return nil, nil
}

func NewMergeItemsTaskExecutor(boardService service.CollectionBoard, collection service.Collection) *MergeItemsTaskExecutor {
	return &MergeItemsTaskExecutor{
		service:           &boardService,
		collectionService: &collection,
	}
}
