package models

import (
	"time"
)

type Collection struct {
	ID             string    `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	CreationDate   time.Time `json:"creationDate,omitempty"`
	ContainedBooks int       `json:"containedBooks,omitempty"`
	OwnerID        string    `json:"ownerID,omitempty"`
	Exclusive      bool      `json:"exclusive,omitempty"`
	ReadCol        bool      `json:"readCol,omitempty"`
	Editable       bool      `json:"editable,omitempty"`
}
