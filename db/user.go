package db

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/TheSgtPepper23/GreenLibrary/models"
	"github.com/TheSgtPepper23/GreenLibrary/services"
)

type UserSQLContext struct {
	conn Database
}

func NewSQLUserContext(pool Database) *UserSQLContext {
	return &UserSQLContext{
		conn: pool,
	}
}

func (c *UserSQLContext) AuthenticateUser(user *models.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var recoveredPassword string
	var id string

	err := c.conn.QueryRow(ctx, `SELECT id, password FROM public.user WHERE email = $1`, strings.ToLower(user.Email)).Scan(&id, &recoveredPassword)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	hashedPass := hashPassword(user.Password)

	if hashedPass != recoveredPassword {
		return "", fmt.Errorf("passwords not match")
	}
	return id, nil
}

func (c *UserSQLContext) UserWizard(email, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	tx, err := c.conn.Begin(ctx)

	if err != nil {
		return err
	}

	hashedPass := hashPassword(password)
	userId := services.GenerateUUID()
	_, err = tx.Exec(ctx, `INSERT INTO public.user (id, email, password) VALUES ($1, $2, $3)`, userId, email, hashedPass)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO public.collection (id, name, creation_date, owner_id, exclusive, read_col, editable) VALUES ($1, $2, $3, $4)`, services.GenerateUUID(), "Leidos", time.Now(), userId, true, true, false)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO public.collection (id, name, creation_date, owner_id, editable) VALUES ($1, $2, $3, $4)`, services.GenerateUUID(), "Por leer", time.Now(), userId, false)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return nil
}

func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}
