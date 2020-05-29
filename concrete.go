package typoless

// concreate
func Count(values ...Field) Call {
	return Call{Prefix: "COUNT", Args: values}
}

var STAR = StringField("*")
