package spoon

// KeyPart holds the columns and the order of arrangement that constitute the key to the index.
type KeyPart struct {
	ColumnName  string
	IsOrderDesc bool
}
