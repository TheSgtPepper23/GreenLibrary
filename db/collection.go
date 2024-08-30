package db

import (
	"context"
	"time"

	"github.com/TheSgtPepper23/GreenLibrary/models"
	"github.com/TheSgtPepper23/GreenLibrary/services"
)

type CollectionSQLContext struct {
	conn Database
}

func NewSQLCollectionContext(pool Database) *CollectionSQLContext {
	return &CollectionSQLContext{
		conn: pool,
	}
}

func (c *CollectionSQLContext) CreateCollection(collection *models.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	collection.ID = services.GenerateUUID()
	collection.CreationDate = time.Now()
	_, err := c.conn.Exec(ctx, `INSERT INTO public.collection (
		id,
		name, 
		creation_date,
		owner_id) 
		VALUES ($1, $2, $3, $4) RETURNING id`, collection.ID, collection.Name, collection.CreationDate, collection.OwnerID)

	if err != nil {
		return err
	}

	return nil
}

func (c *CollectionSQLContext) UpdateCollection(collection *models.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	_, err := c.conn.Exec(ctx, `UPDATE 
		public.collection 
		SET name = $1 
		WHERE id=$2;`, collection.Name, collection.ID)
	return err
}

func (c *CollectionSQLContext) GetCollections() (*[]models.Collection, error) {

	collections := make([]models.Collection, 0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	rows, err := c.conn.Query(ctx, `SELECT 
		c.id, 
		c.name, 
		c.creation_date, 
		c.owner_id,
		COUNT(b.collection_id) as count 
		FROM public.collection c LEFT JOIN public.collection_has_book b 
		on b.collection_id = c.id  
		GROUP BY c.id, c.name 
		ORDER BY c.creation_date desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var collection models.Collection
		err := rows.Scan(&collection.ID, &collection.Name,
			&collection.CreationDate, &collection.OwnerID, &collection.ContainedBooks)
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}

	return &collections, nil
}

func (c *CollectionSQLContext) DeleteCollection(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	tx, err := c.conn.Begin(context.Background())

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM public.collection_has_book WHERE collection_id = $1`, id)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM public.collection WHERE id = $1`, id)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	return nil
}
