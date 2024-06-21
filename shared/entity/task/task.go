package task

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type (
	Type   string
	IDType uuid.UUID
	Status uint
)

const (
	NewBoard   Type = "NewBoard"
	MoveItem   Type = "MoveItem"
	MergeItems Type = "MergeItems"
	ClickItem  Type = "ClickItem"

	Scheduled Status = 0
	Running   Status = 1
	Failed    Status = 2
	Done      Status = 3
)

// Task is a model
type Task struct {
	ID     uint
	UUID   uuid.UUID
	Type   Type
	Status Status

	Args   json.RawMessage
	Result json.RawMessage

	TimeCreated          time.Time
	TimeStartedExecuting time.Time
	TimeDoneExecuting    time.Time
}

// Board is a model
type Board struct {
	ID     uint
	width  uint
	height uint
}
