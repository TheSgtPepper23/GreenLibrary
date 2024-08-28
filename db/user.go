package db

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/TheSgtPepper23/GreenLibrary/models"
)

type UserSQLContext struct {
	conn Database
}

func NewSQLUserContext(pool Database) *UserSQLContext {
	return &UserSQLContext{
		conn: pool,
	}
}

func (c *UserSQLContext) AuthenticateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var recoveredPassword string

	err := c.conn.QueryRow(ctx, `SELECT password FROM public.user WHERE email = $1`, strings.ToLower(user.Email)).Scan(&recoveredPassword)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	hashedPass := hashPassword(user.Password)

	if hashedPass != recoveredPassword {
		return fmt.Errorf("passwords not match")
	}
	return nil
}

func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}
