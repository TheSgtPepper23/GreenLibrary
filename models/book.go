package models

import "time"

type Book struct {
	Id            int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Key           string    `json:"key"`
	AuthorKey     string    `json:"authorKey"`
	ReleaseYear   int       `json:"releaseYear"`
	DateAdded     time.Time `json:"dateAdded"`
	StartReading  time.Time `json:"startReading"`
	FinishReading time.Time `json:"finishReading"`
	CoverURL      string    `json:"coverURL"`
	MyRating      float32   `json:"myRating"`
	AVGRating     float32   `json:"avgRating"`
	Comment       string    `json:"comment"`
	PageCount     int       `json:"pageCount"`
}
