package typoless

import "fmt"

type Literal string

func (v Literal) Name() string {
	return string(v)
}

func Literalf(fmt string, args ...interface{}) LiteralFormat {
	return LiteralFormat{
		Format: fmt,
		Args:   args,
	}
}

type LiteralFormat struct {
	Format string
	Args   []interface{}
}

func (v LiteralFormat) Name() string {
	return fmt.Sprintf(v.Format, v.Args...)
}
func (v LiteralFormat) As(name string) AsField {
	return AsField{Field: v, NewName: name}
}
