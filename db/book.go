package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TheSgtPepper23/GreenLibrary/models"
	"github.com/TheSgtPepper23/GreenLibrary/services"
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

func validateBookReaded(bookID string, userID string, conn *pgxpool.Pool) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	bookReaded := -1

	err := conn.QueryRow(ctx, `SELECT chb.id
		FROM public.collection_has_book chb
		JOIN public.collection c
		ON c.id = chb.collection_id
		WHERE book_id = $1
		AND c.owner_id = $2
		AND finish_reading IS NOT NULL`, bookID, userID).Scan(&bookReaded)

	if err != nil {
		if err != pgx.ErrNoRows {
			return false, err
		} else {
			return false, nil
		}
	}

	return true, nil
}

func validateBookIsStored(bookID string, conn *pgxpool.Pool) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var existingID string

	err := conn.QueryRow(ctx, `SELECT id FROM public.book WHERE "key" = $1`, bookID).Scan(&existingID)
	if err != nil {
		//si el error no es de tipo ErrNoRows significa que algo salió mal
		if err != pgx.ErrNoRows {
			return "", err
		} else {
			return "", nil
		}
	}

	return existingID, nil
}

// Se utiliza para agregar un nuevo libro y asignarlo a una colección
// sa valida que el libro no esté guardado anteriormente para evitar duplicados en la base de datos
func (c *BookSQLContext) CreateNewBook(book *models.Book, userID string) error {
	fmt.Println(userID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	//se revisa si el libro ya existe en la base de datos
	existingID, err := validateBookIsStored(book.Key, c.conn)

	if err != nil {
		return err
	}

	tx, err := c.conn.Begin(ctx)

	if err != nil {
		return err
	}

	//Se asigna el ID. Si no existia se asigna un valor nulo, pero en la parte de creación se va a sustituir. Si ya existe quedó asignado
	book.ID = existingID
	book.DateAdded = time.Now()

	//Si es libro ya se encuentra registrado revisa que no haya sido leido anteriormente
	//Si no está registrado lo registra
	if existingID != "" {
		bookReaded, err := validateBookReaded(existingID, userID, c.conn)

		if err != nil {
			tx.Rollback(ctx)
			return err
		}

		if bookReaded {
			tx.Rollback(ctx)
			return errors.New("book already read")
		}
	} else {
		imgUrl, err := services.ProcessImage(book.CoverURL, book.Key)

		if err != nil {
			fmt.Println("la imagen no se pudo guardar ", err.Error())
			imgUrl = book.CoverURL
		}
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

	//si la fecha de terminado de un libro no es 0 (año 0 o literalmente null) significa que se está marcando como leído
	//por lo tanto se debe de eliminar de las demás colecciones. así es la lógica de la aplicación
	if !book.FinishReading.IsZero() {
		//TODO cambiar esta query para que solo se elimine de las listas del usuario
		_, err = tx.Exec(ctx, `
			DELETE FROM
			public.collection_has_book
			WHERE book_id = $1
		 	AND collection_id IN (SELECT id FROM public.collection WHERE owner_id = $2)`, book.ID, userID)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO public.collection_has_book
			(date_added, book_id, collection_id, rating, comment, finish_reading)
			VALUES($1, $2, $3, $4, $5, $6)`, book.DateAdded, book.ID, book.CollecionID, book.MyRating, book.Comment, book.FinishReading)

		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	} else {
		_, err = tx.Exec(ctx, `INSERT INTO public.collection_has_book
			(date_added, book_id, collection_id)
			VALUES($1, $2, $3)`, book.DateAdded, book.ID, book.CollecionID)

		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return nil
}

func (c *BookSQLContext) UpdateBook(book *models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
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

func (c *BookSQLContext) GetBooksOfCollection(collectionID string) (*[]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	books := make([]models.Book, 0)

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

func (c *BookSQLContext) MoveBook(book *models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := c.conn.Exec(ctx, `UPDATE public.collection_has_book SET collection_id = $1, "comment" = $2, rating = $3 WHERE book_id = $4`, book.CollecionID, book.Comment, book.MyRating, book.ID)
	return err
}

func (c *BookSQLContext) RemoveBookFromCollection(book *models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := c.conn.Exec(ctx, `DELETE FROM public.collection_has_book WHERE book_id = $1 AND collection_id = $2`, book.ID, book.CollecionID)
	return err
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
