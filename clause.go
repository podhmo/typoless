package typoless

import (
	"fmt"
	"strings"
)

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

func From(table tableLike) *FromClause {
	return &FromClause{
		Prefix: "FROM",
		Table:  table,
	}
}

// for FROM

type FromClause struct {
	Prefix string
	Table  tableLike
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
