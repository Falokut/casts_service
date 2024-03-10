package models

type Profession struct {
	ID   int32  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
