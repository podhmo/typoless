package typoless

type tableLike interface {
	TableName() string
}

type Table string
