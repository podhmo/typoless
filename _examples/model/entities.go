package model

type Person struct {
	ID       int64  `db:"id,primarykey,autoincrement"`
	FatherID int64  `db:"father_id"` // todo nullable
	MotherID int64  `db:"mother_id"` // todo nullable
	Name     string `db:"name"`
}
