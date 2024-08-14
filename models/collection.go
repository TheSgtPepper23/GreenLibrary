package models

import (
	"time"
)

type Collection struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	CreationDate   time.Time `json:"creationDate"`
	ContainedBooks int       `json:"containedBooks"`
}
