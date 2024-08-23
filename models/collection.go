package models

import (
	"time"
)

type Collection struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	CreationDate   time.Time `json:"creationDate"`
	ContainedBooks int       `json:"containedBooks"`
	OwnerID        string    `json:"ownerID"`
}
