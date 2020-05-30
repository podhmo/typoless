package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"m/model"
	"m/q"

	"github.com/go-gorp/gorp/v3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

func initDb(ctx context.Context) (*gorp.DbMap, error) {
	rawdb, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	logger := zerolog.New(os.Stdout)
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// populate log pre-fields here before set to
	dbmap := &gorp.DbMap{
		Db: sqldblogger.OpenDriver(
			":memory:",
			rawdb.Driver(),
			zerologadapter.New(logger),
		),
		Dialect: gorp.SqliteDialect{},
	}

	dbmap.AddTableWithName(model.Person{}, "people").SetKeys(true, "id")

	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		return nil, err
	}
	return dbmap, nil
}

func run() error {
	ctx := context.Background()

	fmt.Println("-- setup ----------------------------------------")
	dbmap, err := initDb(ctx)
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

		if err := dbmap.Insert(people[0], people[1], people[2]); err != nil {
			return err
		}
	}

	fmt.Println("-- query all ----------------------------------------")
	{
		var rows []model.Person
		_, err := q.Person.Query(
			q.Select(q.Person.ID, q.Person.Name),
		).DoWithValues(dbmap.Select, &rows)
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
		).Do(dbmap.SelectOne, &ob)
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
		_, err := q.Person.Query(
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
		).DoWithValues(dbmap.Select, &rows)
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
		_, err := q.Person.Query(
			q.Select(
				q.Person.ID,
				q.Person.Name,
				q.Literalf(
					"case when %s=0 AND %s=0 then 1 else 0 end",
					q.Person.MotherID,
					q.Person.FatherID,
				).As("origin"),
			),
		).DoWithValues(dbmap.Select, &rows)
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
