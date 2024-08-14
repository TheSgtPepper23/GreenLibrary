package db

import (
	"context"
	"time"

	"github.com/TheSgtPepper23/GreenLibrary/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookSQLContext struct {
	conn *pgxpool.Pool
}

func NewSQLBookContext(pool *pgxpool.Pool) *BookSQLContext {
	return &BookSQLContext{
		conn: pool,
	}
}

func (c *BookSQLContext) CreateNewBook(book *models.Book) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	var newID int
	err := c.conn.QueryRow(ctx, `INSERT INTO public.book
		( 
		title, 
		author, 
		"key", 
		author_key, 
		release_year, 
		date_added, 
		cover_url,
		avg_rating, 
		page_count, 
		collection_id)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10 ) RETURNING id;`,
		book.Title, book.Author, book.Key, book.AuthorKey, book.ReleaseYear, time.Now(),
		book.CoverURL, book.AVGRating, book.PageCount, book.CollecionID).Scan(newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}
