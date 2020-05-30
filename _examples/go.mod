module m

go 1.14

require (
	github.com/jmoiron/sqlx v1.2.1-0.20190826204134-d7d95172beb5
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/podhmo/typoless v0.0.0
	github.com/rs/zerolog v1.18.0
	github.com/simukti/sqldb-logger v0.0.0-20200401101904-10de8322a496
	google.golang.org/appengine v1.6.6 // indirect
)

replace github.com/podhmo/typoless => ../
