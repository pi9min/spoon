package spoon

import (
	"fmt"
	"strings"
)

// PrimaryKey XXX
type PrimaryKey struct {
	keyParts             []KeyPart
	interleavedTableName string
}

/*// Columns XXX
func (pk *PrimaryKey) Columns() []string {
	cols := make([]string, 0, len(pk.keyParts))
	for _, kp := range pk.keyParts {
		cols = append(cols, kp.ColumnName)
	}

	return cols
}
*/
// ToSQL return primary key sql string
func (pk *PrimaryKey) ToSQL() string {
	var keyPartsStr []string
	for _, kp := range pk.keyParts {
		kps := Quote(kp.ColumnName)
		if kp.IsOrderDesc {
			kps += " DESC"
		}
		keyPartsStr = append(keyPartsStr, kps)
	}

	if pk.interleavedTableName != "" {
		return fmt.Sprintf(`PRIMARY KEY (%s), INTERLEAVE IN PARENT %s`, strings.Join(keyPartsStr, ", "), Quote(pk.interleavedTableName))
	}

	return fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(keyPartsStr, ", "))
}

// AddPrimaryKey XXX
func AddPrimaryKey(keyParts ...KeyPart) *PrimaryKey {
	return &PrimaryKey{
		keyParts:             keyParts,
		interleavedTableName: "",
	}
}

// AddPrimaryKey XXX
func AddPrimaryKeyWithInterleave(interleaveTableName string, keyParts ...KeyPart) *PrimaryKey {
	return &PrimaryKey{
		keyParts:             keyParts,
		interleavedTableName: interleaveTableName,
	}
}
