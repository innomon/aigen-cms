package descriptors

type DataType string

const (
	Int        DataType = "Int"
	Datetime   DataType = "Datetime"
	Text       DataType = "Text"
	String     DataType = "String"
	Lookup     DataType = "Lookup"
	Junction   DataType = "Junction"
	Collection DataType = "Collection"
)

func (d DataType) IsCompound() bool {
	return d == Lookup || d == Junction || d == Collection
}

func (d DataType) IsLocal() bool {
	return d != Junction && d != Collection
}
