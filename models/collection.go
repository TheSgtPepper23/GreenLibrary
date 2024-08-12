package models

import (
	"time"
)

type Collection struct {
	ID             int
	Name           string
	CreationDate   time.Time
	ContainedBooks int
}
