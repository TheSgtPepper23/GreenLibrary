package db

import (
	"context"
	"time"

	"github.com/TheSgtPepper23/GreenLibrary/models"
	"github.com/TheSgtPepper23/GreenLibrary/services"
	"github.com/jackc/pgx/v5"
)

type BookSQLContext struct {
	conn Database
}

func NewSQLBookContext(pool Database) *BookSQLContext {
	return &BookSQLContext{
		conn: pool,
	}
}

func (c *BookSQLContext) CreateNewBook(book *models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	imgUrl, err := services.ProcessImage(book.CoverURL, book.Key)

	//TODO En caso de que no la pueda descargar debe asignarle una portada generica
	if err != nil {
		imgUrl = book.CoverURL
	}

	tx, err := c.conn.Begin(context.Background())

	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	var existingID string
	insertNew := false

	//busca el id del libro que tenga cierta key, si no lo encuentra lanza un error
	err = tx.QueryRow(ctx, `SELECT id FROM public.book WHERE "key" = $1`, book.Key).Scan(&existingID)
	if err != nil {
		//si el error no es de tipo ErrNoRows significa que algo salió mal
		if err != pgx.ErrNoRows {
			return err
		} else {
			//si es del tipo correcto solo signigica que debe de generarse un nuevo libro
			insertNew = true

		}
	}

	//Se asigna el ID. Si no existia se asigna un valor nulo, pero en la parte de creación se va a sustituir. Si ya existe quedó asignado
	book.ID = existingID
	book.DateAdded = time.Now()

	if insertNew {
		book.ID = services.GenerateUUID()
		_, err = tx.Exec(ctx, `INSERT INTO public.book
		(  id, title,  author,  "key",  author_key, 
		release_year,  cover_url, avg_rating,  page_count)
		 VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9 )`,
			book.ID, book.Title, book.Author, book.Key, book.AuthorKey,
			book.ReleaseYear, imgUrl, book.AVGRating, book.PageCount)

		if err != nil {
			return err
		}
	}

	//Se genera la relación del libro (nuevo o existente) con la collección indicada
	_, err = tx.Exec(ctx, `INSERT INTO public.collection_has_book
		(date_added, book_id, collection_id)
		VALUES($1, $2, $3)`, book.DateAdded, book.ID, book.CollecionID)

	if err != nil {
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}

	return nil
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
		cover_url = $6,
		avg_rating = $7,
		page_count = $8,
		WHERE id = $9`,
		book.Title, book.Author, book.Key, book.AuthorKey, book.ReleaseYear,
		book.CoverURL, book.MyRating, book.AVGRating, book.PageCount, book.ID)
	return err
}

func (c *BookSQLContext) GetBookOfCollection(collectionID string) (*[]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	var books []models.Book

	rows, err := c.conn.Query(ctx, `SELECT  b.id, b.title, b.author, b."key", b.author_key, b.release_year,
		chb.date_added, chb.start_reading, chb.finish_reading, b.cover_url,
		chb.rating, chb."comment", b.avg_rating, b.page_count, chb.collection_id
		FROM public.book b left join public.collection_has_book chb on b.id = chb.book_id 
		where chb.collection_id = $1`, collectionID)

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

		//se generan variables temporales para poder escanear los valores de la base de datos que pueden ser nulos
		var dateAdded *time.Time
		var startReading *time.Time
		var finishReading *time.Time
		var myRating *float32
		var avgRating *float32
		var comment *string

		var temp models.Book

		err := rows.Scan(&temp.ID, &temp.Title, &temp.Author, &temp.Key, &temp.AuthorKey,
			&temp.ReleaseYear, &dateAdded, &startReading, &finishReading,
			&temp.CoverURL, &myRating, &comment, &avgRating, &temp.PageCount,
			&temp.CollecionID)
		if err != nil {
			return err
		}

		//se asignan los valores temporales a los valores finales en caso de no ser nulos. De lo contrario los valores se dejan en cero
		if dateAdded != nil {
			temp.DateAdded = *dateAdded
		}

		if startReading != nil {
			temp.StartReading = *startReading
		}

		if finishReading != nil {
			temp.FinishReading = *finishReading
		}

		if myRating != nil {
			temp.MyRating = *myRating
		}

		if avgRating != nil {
			temp.AVGRating = *avgRating
		}

		if comment != nil {
			temp.Comment = *comment
		}

		*target = append(*target, temp)
	}

	return nil
}
