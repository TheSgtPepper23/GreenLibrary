package models

type MainResponse struct {
	Collections   *[]Collection
	Books         *[]Book
	NewBook       *Book
	NewCollection *Collection
}
