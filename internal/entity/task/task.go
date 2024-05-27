package task

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Type string

const NewBoard Type = "NewBoard"

type Status uint

const Scheduled Status = 0
const Running Status = 1
const Failed Status = 2
const Done Status = 3

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
