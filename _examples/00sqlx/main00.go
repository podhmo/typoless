package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"m/model"
	"m/q"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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
		people := []*model.Person{
			&model.Person{
				ID:       1,
				FatherID: 2,
				MotherID: 3,
				Name:     "me",
			},
			&model.Person{
				ID:   2,
				Name: "F",
			},
			&model.Person{
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
		var rows []model.Person
		err := q.Person.Query(
			q.Select(q.Person.ID, q.Person.Name),
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
		var ob model.Person
		err := q.Person.Query(
			q.Where(q.Person.Name.Compare("= ?", "me")),
			q.Select(q.Person.ID, q.Person.Name),
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

		p := q.Person.As("p")
		father := q.Person.As("father")
		mother := q.Person.As("mother")

		var rows []view
		err := q.Person.Query(
			q.From(
				p.
					Join(father, q.On(p.FatherID, father.ID)).
					Join(mother, q.On(p.MotherID, mother.ID)),
			),
			q.Select(
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

	fmt.Println("-- with literalf ----------------------------------------")
	{
		type result struct {
			ID     int64  `db:"id"`
			Name   string `db:"name"`
			Origin bool   `db:"origin"`
		}
		var rows []result
		err := q.Person.Query(
			q.Select(
				q.Person.ID,
				q.Person.Name,
				q.Literalf(
					"case when %s then 1 else 0 end",
					q.And(
						q.Person.MotherID.Compare("=", 0),
						q.Person.FatherID.Compare("=", 0),
					),
				).As("origin"),
			),
			q.OrderBy(q.Desc(q.Person.ID), q.Asc(q.Person.Name)),
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

var _schema = `
CREATE TABLE people (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  father_id INTEGER,
  mother_id INTEGER,
  name TEXT
);
`
