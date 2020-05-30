package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/podhmo/typoless"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

func initDb(ctx context.Context) (*sqlx.DB, error) {
	rawdb, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	logger := zerolog.New(os.Stdout)
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// populate log pre-fields here before set to
	db := sqlx.NewDb(sqldblogger.OpenDriver(
		":memory:",
		rawdb.Driver(),
		zerologadapter.New(logger),
	), "sqlite3")
	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.MustExecContext(ctx, _schema)
	return db, nil
}

func run() error {
	ctx := context.Background()

	fmt.Println("-- setup ----------------------------------------")
	db, err := initDb(ctx)
	if err != nil {
		return err
	}

	fmt.Println("-- insert ----------------------------------------")
	{
		people := []*Person{
			&Person{
				ID:       1,
				FatherID: 2,
				MotherID: 3,
				Name:     "me",
			},
			&Person{
				ID:   2,
				Name: "F",
			},
			&Person{
				ID:   3,
				Name: "M",
			},
		}

		tx := db.MustBegin()
		// TODO
		if _, err := tx.NamedExec(`
INSERT INTO people (id, name, father_id, mother_id) VALUES (:id, :name, :father_id, :mother_id)
`, people); err != nil {
			return err
		}
		tx.Commit()
	}

	fmt.Println("-- query all ----------------------------------------")
	{
		var rows []Person
		err := PersonD.Query(
			Select(PersonD.ID, PersonD.Name),
		).Do(db.Select, &rows)
		if err != nil {
			return err
		}
		log.Println("All rows:")
		for i, ob := range rows {
			log.Printf("    %d: %#+v\n", i, ob)
		}
	}

	fmt.Println("-- query one ----------------------------------------")
	{
		// TODO: limit
		var ob Person
		err := PersonD.Query(
			Where(PersonD.Name.Compare("= ?", "me")),
			Select(PersonD.ID, PersonD.Name),
		).Do(db.Get, &ob)
		if err != nil {
			return err
		}
		log.Printf("\tgot, %#v\n", ob)
	}

	fmt.Println("-- join query ----------------------------------------")
	{
		type view struct {
			ID         int64  `db:"id"`
			Name       string `db:"name"`
			FatherName string `db:"father_name"`
			MotherName string `db:"mother_name"`
		}

		p := PersonD.As("p")
		father := PersonD.As("father")
		mother := PersonD.As("mother")

		var rows []view
		err := PersonD.Query(
			From(
				p.
					Join(father, On(p.FatherID, father.ID)).
					Join(mother, On(p.MotherID, mother.ID)),
			),
			Select(
				p.ID,
				p.Name,
				father.Name.As("father_name"),
				mother.Name.As("mother_name"),
			),
		).Do(db.Select, &rows)
		if err != nil {
			return err
		}
		log.Println("All rows:")
		for i, ob := range rows {
			log.Printf("    %d: %#+v\n", i, ob)
		}
	}
	return nil
}

func main() {
	log.SetPrefix("*********** ")
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatalf("!!+%v", err)
	}
}

type Person struct {
	ID       int64  `db:"id"`
	FatherID int64  `db:"father_id"` // todo nullable
	MotherID int64  `db:"mother_id"` // todo nullable
	Name     string `db:"name"`
}

var _schema = `
CREATE TABLE people (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  father_id INTEGER,
  mother_id INTEGER,
  name TEXT
);
`

type PersonDefinition struct {
	typoless.Table
	ID       typoless.Int64Field
	FatherID typoless.Int64Field
	MotherID typoless.Int64Field
	Name     typoless.StringField
}

func (d *PersonDefinition) As(name string) *PersonDefinition {
	new := *d
	typoless.Alias(&new, d, name)
	return &new
}

var PersonD = PersonDefinition{
	Table:    typoless.Table("people"),
	ID:       typoless.Int64Field("id"),
	FatherID: typoless.Int64Field("father_id"),
	MotherID: typoless.Int64Field("mother_id"),
	Name:     typoless.StringField("name"),
}

var (
	Select = PersonD.Select
	From   = PersonD.From
	Where  = PersonD.Where
	On     = typoless.On
)
