package db

import (
	"context"
	"time"

	"github.com/TheSgtPepper23/GreenLibrary/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CollectionSQLContext struct {
	conn *pgxpool.Pool
}

func NewSQLCollectionContext(pool *pgxpool.Pool) *CollectionSQLContext {
	return &CollectionSQLContext{
		conn: pool,
	}
}

func (c *CollectionSQLContext) CreateCollection(collection *models.Collection) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	var newID int
	err := c.conn.QueryRow(ctx, `INSERT INTO public.collection (
		name, 
		creation_date) 
		VALUES ($1, $2) RETURNING id`, collection.Name, collection.CreationDate).Scan(newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (c *CollectionSQLContext) UpdateCollection(collection *models.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	_, err := c.conn.Exec(ctx, `UPDATE 
		public.collection 
		SET name = $1, 
		creation_date= $2 
		WHERE id=$3;`, collection.Name, collection.CreationDate, collection.ID)
	return err
}

func (c *CollectionSQLContext) GetCollections() (*[]models.Collection, error) {
	var collections []models.Collection

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	rows, err := c.conn.Query(ctx, `SELECT 
		c.id, 
		c.name, 
		c.creation_date, 
		COUNT(b.id) as count 
		FROM public.collection c LEFT JOIN public.book b 
		on b.id = c.id  
		GROUP BY c.id, c.name 
		ORDER BY c.creation_date desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var collection models.Collection
		err := rows.Scan(&collection.ID, &collection.Name,
			&collection.CreationDate, &collection.ContainedBooks)
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}

	return &collections, nil
}
