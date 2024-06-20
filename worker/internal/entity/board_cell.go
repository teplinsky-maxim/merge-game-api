package entity

import (
	"time"
)

// BoardCell is a model
type BoardCell struct {
	ID               uint
	BoardID          uint
	CellW            uint
	CellH            uint
	CollectionID     uint
	CollectionItemID uint

	TimeCreated time.Time
}
