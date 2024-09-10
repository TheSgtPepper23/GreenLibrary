package db

import (
	"context"
	"errors"
	"time"

	"github.com/TheSgtPepper23/GreenLibrary/models"
	"github.com/TheSgtPepper23/GreenLibrary/services"
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

func (c *CollectionSQLContext) CreateCollection(collection *models.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := c.conn.Exec(ctx, `INSERT INTO public.collection (
		id,
		name,
		creation_date,
		owner_id,
		exclusive)
		VALUES ($1, $2, $3, $4, $5)`, services.GenerateUUID(), collection.Name, time.Now(), collection.OwnerID, collection.Exclusive)

	if err != nil {
		return err
	}

	return nil
}

func (c *CollectionSQLContext) UpdateCollection(collection *models.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := c.conn.Exec(ctx, `UPDATE
		public.collection
		SET name = $1,
		exclusive = $2
		WHERE id=$3;`, collection.Name, collection.Exclusive, collection.ID)
	return err
}

func (c *CollectionSQLContext) GetCollections(ownerID string) (*[]models.Collection, error) {

	collections := make([]models.Collection, 0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := c.conn.Query(ctx, `SELECT
		c.id,
		c.name,
		c.creation_date,
		c.owner_id,
		c.exclusive,
		c.read_col,
		c.editable,
		COUNT(b.collection_id) as count
		FROM public.collection c LEFT JOIN public.collection_has_book b
		on b.collection_id = c.id
		WHERE c.owner_id = $1
		GROUP BY c.id, c.name
		ORDER BY c.creation_date desc`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var collection models.Collection
		err := rows.Scan(&collection.ID, &collection.Name,
			&collection.CreationDate, &collection.OwnerID,
			&collection.Exclusive, &collection.ReadCol,
			&collection.Editable, &collection.ContainedBooks)
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}

	return &collections, nil
}

func (c *CollectionSQLContext) DeleteCollection(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	tx, err := c.conn.Begin(ctx)
	if err != nil {
		return err
	}

	idEditable := true
	err = tx.QueryRow(ctx, `SELECT editable FROM public.collection WHERE id = $1`, id).Scan(&idEditable)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if !idEditable {
		tx.Rollback(ctx)
		return errors.New("collection is not editable")
	}

	_, err = tx.Exec(ctx, `DELETE FROM public.collection_has_book WHERE collection_id = $1`, id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM public.collection WHERE id = $1`, id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return nil
}
