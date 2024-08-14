package models

import "os"

type ServerContext struct {
	AllowCreation bool
	Data          MainResponse
}

func NewContext() *ServerContext {
	allowCreate := os.Getenv("ALLOW_CREATION")
	return &ServerContext{
		AllowCreation: allowCreate == "1",
		Data: MainResponse{
			Collections: &[]Collection{},
			Books:       &[]Book{},
		},
	}
}
