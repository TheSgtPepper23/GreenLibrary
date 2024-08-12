package services

import (
	"testing"
)

func TestSQLCreateCollection(t *testing.T) {
	_, err := SearchBook("el señor de los anillos")

	if err != nil {
		t.Log(err.Error())
		t.Error(err)
	}
}
