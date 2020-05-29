package typoless

import (
	"fmt"
)

type Query struct {
	SelectClause *SelectClause
	FromClause   *FromClause
	WhereClause  *WhereClause
}

func (q *Query) Do(
	fn func(ob interface{}, stmt string, args ...interface{}) error,
	ob interface{},
) error {
	return fn(
		ob,
		fmt.Sprintf("%v %v %v", q.SelectClause, q.FromClause, q.WhereClause),
		Values(q.WhereClause.Value)...,
	)
}
func (q *Query) DoWithValue(
	fn func(ob interface{}, stmt string, args ...interface{}) (interface{}, error),
	ob interface{},
) (interface{}, error) {
	return fn(
		ob,
		fmt.Sprintf("%v %v %v", q.SelectClause, q.FromClause, q.WhereClause),
		Values(q.WhereClause.Value)...,
	)
}

func (q *Query) DoWithValues(
	fn func(ob interface{}, stmt string, args ...interface{}) ([]interface{}, error),
	ob interface{},
) ([]interface{}, error) {
	return fn(
		ob,
		fmt.Sprintf("%v %v %v", q.SelectClause, q.FromClause, q.WhereClause),
		Values(q.WhereClause.Value)...,
	)
}
