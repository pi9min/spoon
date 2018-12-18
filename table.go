package spoon

import (
	"fmt"
	"strings"
)

// Table is mapping struct info
type Table struct {
	name       string
	columns    []*Column
	primaryKey *PrimaryKey
	indexes    Indexes
}

func newTable(name string, columns []*Column, pk *PrimaryKey, indexes Indexes) *Table {
	return &Table{
		name:       name,
		columns:    columns,
		primaryKey: pk,
		indexes:    indexes,
	}
}

func (t *Table) Indexes() Indexes {
	return t.indexes
}

func (t *Table) CreateTableSchema() string {
	ss := make([]string, 0, len(t.columns)+2)
	ss = append(ss, fmt.Sprintf("CREATE TABLE %s (", Quote(t.name)))
	for i := range t.columns {
		c := t.columns[i]
		ss = append(ss, fmt.Sprintf("    %s,", c.ToSQL()))
	}
	ss = append(ss, fmt.Sprintf(") %s", t.primaryKey.ToSQL()))
	return strings.Join(ss, "\n")
}

func (t *Table) DropTableSchema() string {
	return fmt.Sprintf("DROP TABLE %s", Quote(t.name))
}
