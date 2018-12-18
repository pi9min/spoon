package spoon

import (
	"fmt"
	"strings"
)

// Indexes are alias of index slices.
type Indexes []*Index

// Index holds the necessary information to construct Index.
type Index struct {
	name         string
	tableName    string
	isUnique     bool
	nullFiltered bool
	keyParts     []KeyPart
}

// CreateIndexSchema return `CREATE INDEX` schema.
func (i *Index) CreateIndexSchema() string {
	var keyPartsStr []string

	for _, kp := range i.keyParts {
		kps := Quote(kp.ColumnName)
		if kp.IsOrderDesc {
			kps += " DESC"
		}
		keyPartsStr = append(keyPartsStr, kps)
	}

	words := make([]string, 0, 4)
	words = append(words, "CREATE")
	if i.isUnique {
		words = append(words, "UNIQUE")
	}
	if i.nullFiltered {
		words = append(words, "NULL_FILTERED")
	}
	words = append(words, "INDEX")

	schema := fmt.Sprintf(
		"%s %s ON %s (%s)",
		strings.Join(words, " "),
		Quote(i.name),
		Quote(i.tableName),
		strings.Join(keyPartsStr, ", "),
	)

	return schema
}

// DropIndexSchema return `DROP INDEX` schema.
func (i *Index) DropIndexSchema() string {
	schema := fmt.Sprintf("DROP INDEX %s", Quote(i.name))
	return schema
}

// AddIndex creates Index.
func AddIndex(idxName, tableName string, nullFiltered bool, keyParts ...KeyPart) *Index {
	return &Index{
		name:         idxName,
		keyParts:     keyParts,
		tableName:    tableName,
		isUnique:     false,
		nullFiltered: nullFiltered,
	}
}

// AddUniqueIndex creates Unique Index.
func AddUniqueIndex(idxName, tableName string, nullFiltered bool, keyParts ...KeyPart) *Index {
	return &Index{
		name:         idxName,
		keyParts:     keyParts,
		tableName:    tableName,
		isUnique:     true,
		nullFiltered: nullFiltered,
	}
}
