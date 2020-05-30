package typoless

import (
	"fmt"
	"strings"
)

// https://www.sqlite.org/syntaxdiagrams.html#select-stmt

// for SELECT

func Select(fields ...Field) *SelectClause {
	return &SelectClause{Prefix: "SELECT", Fields: fields}
}

type SelectClause struct {
	Prefix string
	Fields []Field
}

func (q *SelectClause) String() string {
	names := make([]string, len(q.Fields))
	for i, f := range q.Fields {
		names[i] = f.Name()
	}
	return fmt.Sprintf("%s %s", q.Prefix, strings.Join(names, ", "))
}

func From(table tablelike) *FromClause {
	return &FromClause{
		Prefix: "FROM",
		Table:  table,
	}
}

// for FROM

type FromClause struct {
	Prefix string
	Table  tablelike
}

func (q *FromClause) String() string {
	return fmt.Sprintf("%s %s", q.Prefix, q.Table.TableName())
}

func Where(ops ...op) *WhereClause {
	values := make([]interface{}, len(ops))
	for i, op := range ops {
		values[i] = op
	}

	value := And(values...)
	value.WithoutParen = true
	return &WhereClause{Prefix: "WHERE", Value: value}
}

// for WHERE

type WhereClause struct {
	Prefix string
	Value  *Mop
}

func (q *WhereClause) String() string {
	if len(q.Value.Values) == 0 {
		return ""
	}

	value := Replace(q.Value, "")
	return fmt.Sprintf("%s %v", q.Prefix, value)
}

// for Order by
type OrderByClause struct {
	Prefix string
	Args   []Ordering
}

func OrderBy(args ...Ordering) *OrderByClause {
	return &OrderByClause{
		Prefix: "ORDER BY",
		Args:   args,
	}
}
func (q *OrderByClause) String() string {
	if q == nil {
		return ""
	}
	buf := make([]string, len(q.Args))
	for i := 0; i < len(q.Args); i++ {
		buf[i] = q.Args[i].Name()
	}
	return q.Prefix + " " + strings.Join(buf, ", ")
}
