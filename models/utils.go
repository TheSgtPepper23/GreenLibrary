package models

type OrderOption int

const (
	NameAsc OrderOption = iota
	DateAsc
	NameDesc
	DateDesc
)
