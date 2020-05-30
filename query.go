package typoless

import (
	"strings"
)

type Query struct {
	SelectClause *SelectClause
	FromClause   *FromClause
	WhereClause  *WhereClause

	// GroupByClause *GroupByClause
	OrderByClause *OrderByClause
	// LimitClause   *LimitClause
	// OffsetClause  *OffsetClause
}

func (q *Query) Stmt() string {
	buf := make([]string, 0, 5)
	buf = append(buf, q.SelectClause.String())
	buf = append(buf, q.FromClause.String())
	buf = append(buf, q.WhereClause.String())
	if q.OrderByClause != nil {
		buf = append(buf, q.OrderByClause.String())
	}
	return strings.Join(buf, " ")
}

func (q *Query) Do(
	fn func(ob interface{}, stmt string, args ...interface{}) error,
	ob interface{},
) error {
	return fn(ob, q.Stmt(), Values(q.WhereClause.Value)...)
}
func (q *Query) DoWithValue(
	fn func(ob interface{}, stmt string, args ...interface{}) (interface{}, error),
	ob interface{},
) (interface{}, error) {
	return fn(ob, q.Stmt(), Values(q.WhereClause.Value)...)
}
func (q *Query) DoWithValues(
	fn func(ob interface{}, stmt string, args ...interface{}) ([]interface{}, error),
	ob interface{},
) ([]interface{}, error) {
	return fn(ob, q.Stmt(), Values(q.WhereClause.Value)...)
}
