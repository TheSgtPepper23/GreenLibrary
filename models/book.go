package models

import "time"

type Book struct {
	Id            int
	Title         string
	Author        string
	Key           string
	AuthorKey     string
	ReleaseYear   int
	DateAdded     time.Time
	StartReading  time.Time
	FinishReading time.Time
	CoverURL      string
	MyRating      float32
	AVGRating     float32
	Comment       string
	PageCount     int
}
