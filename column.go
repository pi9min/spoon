package spoon

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	maxStringLength = 2621440
	maxByteLength   = 10485760
)

// Column is mapping struct field value.
type Column struct {
	name        string
	isNull      bool
	size        int
	reflectType reflect.Type
}

func newColumn(name string, tags map[string]string, rt reflect.Type) *Column {
	isNull, size, err := parseTags(tags)
	if err != nil {
		panic(err)
	}

	return &Column{
		name:        name,
		isNull:      isNull,
		size:        size,
		reflectType: rt,
	}
}

// ToSQL is convert struct value to sql.
// ToSQL convert spanner type from reflect.Type and size
func (c *Column) ToSQL() string {
	tStr, tNull := parseTypeToString(c.reflectType, c.size)
	// Always NOT NULL if both nulls are not satisfied
	if !(c.isNull || tNull) {
		tStr += " NOT NULL"
	}
	return fmt.Sprintf("%s %s", Quote(c.name), tStr)
}

func parseTypeToString(t reflect.Type, size int) (string, bool) {
	switch t.Kind() {
	// Recursive
	case reflect.Ptr:
		return parseTypeToString(t.Elem(), size)
	case reflect.Bool:
		return "BOOL", false
	case reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16, reflect.Int, reflect.Uint, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64:
		return "INT64", false
	case reflect.Float32, reflect.Float64:
		return "FLOAT64", false
	case reflect.Slice:
		switch t.Elem().Kind() {
		case reflect.Uint8: // []byte
			if size < 1 || maxByteLength < size {
				return "BYTES(MAX)", false // MAX=10485760(10MiB)
			}
			return fmt.Sprintf("BYTES(%d)", size), false
		default:
			typeStr, isNull := parseTypeToString(t.Elem(), size)
			return array(typeStr), isNull
		}
	}

	switch t.Name() {
	case "Time":
		return "TIMESTAMP", false
	case "NullBool": // https://godoc.org/cloud.google.com/go/spanner#NullBool
		return "BOOL", true
	case "NullDate", "Date": // https://godoc.org/cloud.google.com/go/spanner#NullDate, https://godoc.org/cloud.google.com/go/civil#Date
		return "DATE", true
	case "NullFloat64": // https://godoc.org/cloud.google.com/go/spanner#NullFloat64
		return "FLOAT64", true
	case "NullInt64": // https://godoc.org/cloud.google.com/go/spanner#NullInt64
		return "INT64", true
	case "NullTime": // https://godoc.org/cloud.google.com/go/spanner#NullTime
		return "TIMESTAMP", true
	}

	// Process the following as a character string.
	var isNull bool
	// https://godoc.org/cloud.google.com/go/spanner#NullString
	if t.Name() == "NullString" {
		isNull = true
	}

	if size < 1 || maxStringLength < size {
		return "STRING(MAX)", isNull // MAX=2621440(2.5mebichars)
	}

	return fmt.Sprintf("STRING(%d)", size), isNull
}

func array(s string) string {
	return "ARRAY<" + s + ">"
}

func parseTags(tags map[string]string) (bool, int, error) {
	var isNull bool

	if _, ok := tags["nullable"]; ok {
		isNull = true
	}

	// size tag
	sizeStr, ok := tags["size"]
	if !ok {
		return isNull, 0, nil
	}

	s, err := strconv.Atoi(sizeStr)
	if err != nil {
		return false, 0, err
	}

	return isNull, s, nil

}
