package q

import "github.com/podhmo/typoless"

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

var Person = PersonDefinition{
	Table:    typoless.Table("people"),
	ID:       typoless.Int64Field("id"),
	FatherID: typoless.Int64Field("father_id"),
	MotherID: typoless.Int64Field("mother_id"),
	Name:     typoless.StringField("name"),
}

var (
	Select  = Person.Select
	From    = Person.From
	Where   = Person.Where
	OrderBy = Person.OrderBy
)

var (
	On   = typoless.On
	Asc  = typoless.Asc
	Desc = typoless.Desc
)

var (
	Literalf = typoless.Literalf
	And      = typoless.And
	Or       = typoless.Or
)
