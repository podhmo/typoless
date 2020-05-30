package typoless

// ops
func And(values ...interface{}) *Mop {
	return &Mop{Op: "AND", Values: values}
}
func Or(values ...interface{}) *Mop {
	return &Mop{Op: "OR", Values: values}
}
func Not(value interface{}) *Uop {
	return &Uop{Op: "NOT", Value: value}
}

// for order by
type Ordering = *LiteralFormat

func Desc(field Field) Ordering {
	return &LiteralFormat{Format: "%s DESC", Args: []interface{}{field.Name()}}
}
func Asc(field Field) Ordering {
	return &LiteralFormat{Format: "%s ASC", Args: []interface{}{field.Name()}}
}

// fields
func Count(values ...Field) Call {
	return Call{Prefix: "COUNT", Args: values}
}

var STAR = StringField("*")
