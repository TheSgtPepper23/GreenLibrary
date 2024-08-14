package db

import (
	"context"
	"time"

	"github.com/TheSgtPepper23/GreenLibrary/models"
	"github.com/jackc/pgx/v5"
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

func (c *BookSQLContext) UpdateBook(book *models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	_, err := c.conn.Exec(ctx, `UPDATE public.book SET 
		title = $1,
		author = $2,
		"key" = $3,
		author_key = $4,
		release_year = $5,
		start_reading = $6,
		finish_reading = $7,
		cover_url = $8,
		my_rating = $9,
		"comment" = $10,
		avg_rating = $11,
		page_count = $12,
		collection_id = $13
		WHERE id = $14`,
		book.Title, book.Author, book.Key, book.AuthorKey, book.ReleaseYear,
		book.StartReading, book.FinishReading, book.CoverURL, book.MyRating,
		book.Comment, book.AVGRating, book.PageCount, book.CollecionID, book.ID)
	return err
}

func (c *BookSQLContext) GetBookOfCollection(collectionID int) (*[]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	var books []models.Book

	rows, err := c.conn.Query(ctx, `SELECT 
		id, title, author, "key", author_key, release_year,
		date_added, start_reading, finish_reading, cover_url,
		my_rating, "comment", avg_rating, page_count, collection_id
		FROM public.book 
		WHERE collection_id = $1
	`, collectionID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	err = scanBooks(rows, &books)
	if err != nil {
		return nil, err
	}

	return &books, nil
}

func (c *BookSQLContext) GetBooks() (*[]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	var books []models.Book

	rows, err := c.conn.Query(ctx, `SELECT 
		id, title, author, "key", author_key, release_year,
		date_added, start_reading, finish_reading, cover_url,
		my_rating, "comment", avg_rating, page_count, collection_id
		FROM public.book`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	err = scanBooks(rows, &books)
	if err != nil {
		return nil, err
	}

	return &books, nil
}

func scanBooks(rows pgx.Rows, target *[]models.Book) error {
	for rows.Next() {
		var temp models.Book
		err := rows.Scan(&temp.ID, &temp.Title, &temp.Author, &temp.Key, &temp.AuthorKey,
			&temp.ReleaseYear, &temp.DateAdded, &temp.StartReading, &temp.FinishReading,
			&temp.CoverURL, &temp.MyRating, &temp.Comment, &temp.AVGRating, &temp.PageCount,
			&temp.CollecionID)
		if err != nil {
			return err
		}
		*target = append(*target, temp)
	}

	return nil
}
