package models

import "time"

type Book struct {
	ID            string    `json:"id,omitempty"`
	Title         string    `json:"title,omitempty"`
	Author        string    `json:"author,omitempty"`
	Key           string    `json:"key,omitempty"`
	AuthorKey     string    `json:"authorKey,omitempty"`
	ReleaseYear   int       `json:"releaseYear,omitempty"`
	DateAdded     time.Time `json:"dateAdded,omitempty"`
	StartReading  time.Time `json:"startReading,omitempty"`
	FinishReading time.Time `json:"finishReading,omitempty"`
	CoverURL      string    `json:"coverURL,omitempty"`
	MyRating      float32   `json:"myRating,omitempty"`
	AVGRating     float32   `json:"avgRating,omitempty"`
	Comment       string    `json:"comment,omitempty"`
	PageCount     int       `json:"pageCount,omitempty"`
	CollecionID   string    `json:"collectionID,omitempty"`
}
