package typoless

import (
	"fmt"
	"strings"
)

type Field interface {
	Name() string
}

type AsField struct {
	NewName string
	Field   Field
}

func (as AsField) Name() string {
	return fmt.Sprintf("%s as %s", as.Field.Name(), as.NewName)
}

type Int64Field string

func (f Int64Field) Name() string {
	return string(f)
}
func (f Int64Field) As(name string) AsField {
	return AsField{NewName: name, Field: f}
}
func (f Int64Field) Compare(op string, value int64) *Bop {
	return &Bop{
		Op:           op, // e.g. "= ?"
		Left:         f.Name(),
		Right:        value,
		WithoutParen: true,
	}
}

type StringField string

func (f StringField) Name() string {
	return string(f)
}
func (f StringField) As(name string) AsField {
	return AsField{NewName: name, Field: f}
}
func (f StringField) Compare(op string, value string) *Bop {
	return &Bop{
		Op:           op, // e.g. "= ?"
		Left:         f.Name(),
		Right:        value,
		WithoutParen: true,
	}
}

type Call struct {
	Prefix string
	Args   []Field
}

func (c Call) Name() string {
	names := make([]string, len(c.Args))
	for i, f := range c.Args {
		names[i] = f.Name()
	}
	return fmt.Sprintf("%s(%s)", c.Prefix, strings.Join(names, ", "))
}
