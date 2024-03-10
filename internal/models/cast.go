package models

type Cast struct {
	Persons []Person `db:"persons" json:"persons"`
}
