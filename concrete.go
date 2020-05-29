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

// fields
func Count(values ...Field) Call {
	return Call{Prefix: "COUNT", Args: values}
}

var STAR = StringField("*")
