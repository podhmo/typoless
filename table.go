package typoless

import "fmt"

type tablelike interface {
	TableName() string
}

type Table string

func (t Table) TableName() string {
	return string(t)
}
func (t Table) As(name string) Table {
	return Table(fmt.Sprintf("%s as %s", t, name))
}
func (t Table) Join(rhs tablelike, on string) JoinedTable {
	return JoinedTable{lhs: t, rhs: rhs, on: on, joiner: "JOIN"}
}
func (t Table) LeftOuterJoin(rhs tablelike, on string) JoinedTable {
	return JoinedTable{lhs: t, rhs: rhs, on: on, joiner: "LEFT OUTER JOIN"}
}
func (t Table) RightOuterJoin(rhs tablelike, on string) JoinedTable {
	return JoinedTable{lhs: t, rhs: rhs, on: on, joiner: "RIGHT OUTER JOIN"}
}
func (t Table) FullOuterJoin(rhs tablelike, on string) JoinedTable {
	return JoinedTable{lhs: t, rhs: rhs, on: on, joiner: "FULL OUTER JOIN"}
}

type JoinedTable struct {
	lhs    tablelike
	rhs    tablelike
	joiner string
	on     string
}

func (t JoinedTable) TableName() string {
	return fmt.Sprintf("%s %s %s %s",
		t.lhs.TableName(),
		t.joiner,
		t.rhs.TableName(),
		t.on,
	)
}
func (t JoinedTable) Join(rhs tablelike, on string) JoinedTable {
	return JoinedTable{lhs: t, rhs: rhs, on: on, joiner: "JOIN"}
}
func (t JoinedTable) LeftOuterJoin(rhs tablelike, on string) JoinedTable {
	return JoinedTable{lhs: t, rhs: rhs, on: on, joiner: "LEFT OUTER JOIN"}
}
func (t JoinedTable) RightOuterJoin(rhs tablelike, on string) JoinedTable {
	return JoinedTable{lhs: t, rhs: rhs, on: on, joiner: "RIGHT OUTER JOIN"}
}
func (t JoinedTable) FullOuterJoin(rhs tablelike, on string) JoinedTable {
	return JoinedTable{lhs: t, rhs: rhs, on: on, joiner: "FULL OUTER JOIN"}
}

func On(lhs, rhs Field) string {
	return fmt.Sprintf("ON %s=%s", lhs.Name(), rhs.Name())
}

func (t Table) Query(
	options ...func(*Query),
) *Query {
	q := &Query{
		FromClause:   From(t),
		SelectClause: Select(STAR),
		WhereClause:  Where(),
	}
	for _, opt := range options {
		opt(q)
	}
	return q
}

func (Table) Select(fields ...Field) func(*Query) {
	return func(q *Query) {
		q.SelectClause = Select(fields...)
	}
}

func (Table) Where(ops ...op) func(*Query) {
	return func(q *Query) {
		q.WhereClause = Where(ops...)
	}
}

func (Table) From(t tablelike) func(*Query) {
	return func(q *Query) {
		q.FromClause = From(t)
	}
}
