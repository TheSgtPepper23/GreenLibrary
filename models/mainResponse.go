package models

type MainResponse struct {
	Collections   *[]Collection
	Books         *[]Book
	NewBook       *Book
	NewCollection *Collection
}

func (r *MainResponse) ChangeBooks(books *[]Book) {
	r.Books = books
}

func (r *MainResponse) ChangeCollections(collections *[]Collection) {
	r.Collections = collections
}

func (r *MainResponse) SetNewCollections(collection *Collection) {
	r.NewCollection = collection
}

func (r *MainResponse) SetNewBook(book *Book) {
	r.NewBook = book
}
