package spoon

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"cloud.google.com/go/spanner"
)

func TestColumn_ParseTypeToString(t *testing.T) {
	tests := []struct {
		name          string
		inputType     interface{}
		size          int
		expectTypeStr string
		expectIsNull  bool
	}{
		{name: "bool", inputType: true, size: 0, expectTypeStr: "BOOL", expectIsNull: false},
		{name: "*bool pointer case", inputType: (*bool)(nil), size: 0, expectTypeStr: "BOOL", expectIsNull: false},
		{name: "int8", inputType: int8(0), size: 0, expectTypeStr: "INT64", expectIsNull: false},
		{name: "int16", inputType: int16(0), size: 0, expectTypeStr: "INT64", expectIsNull: false},
		{name: "int32", inputType: int32(0), size: 0, expectTypeStr: "INT64", expectIsNull: false},
		{name: "int", inputType: int(0), size: 0, expectTypeStr: "INT64", expectIsNull: false},
		{name: "int64", inputType: int64(0), size: 0, expectTypeStr: "INT64", expectIsNull: false},
		{name: "spanner.NullInt64", inputType: spanner.NullInt64{}, size: 0, expectTypeStr: "INT64", expectIsNull: true},
		{name: "uint8", inputType: uint8(0), size: 0, expectTypeStr: "INT64", expectIsNull: false},
		{name: "uint16", inputType: uint16(0), size: 0, expectTypeStr: "INT64", expectIsNull: false},
		{name: "uint32", inputType: uint32(0), size: 0, expectTypeStr: "INT64", expectIsNull: false},
		{name: "uint", inputType: uint(0), size: 0, expectTypeStr: "INT64", expectIsNull: false},
		{name: "uint64", inputType: uint64(0), size: 0, expectTypeStr: "INT64", expectIsNull: false},
		{name: "float32", inputType: float32(0), size: 0, expectTypeStr: "FLOAT64", expectIsNull: false},
		{name: "float64", inputType: float64(0), size: 0, expectTypeStr: "FLOAT64", expectIsNull: false},
		{name: "spanner.NullFloat64", inputType: spanner.NullFloat64{}, size: 0, expectTypeStr: "FLOAT64", expectIsNull: true},
		{name: "*spanner.NullFloat64", inputType: &spanner.NullFloat64{}, size: 0, expectTypeStr: "FLOAT64", expectIsNull: true},
		{name: "[]byte size:0", inputType: []byte{}, size: 0, expectTypeStr: "BYTES(MAX)", expectIsNull: false},
		{name: "[]byte size:1", inputType: []byte{}, size: 1, expectTypeStr: "BYTES(1)", expectIsNull: false},
		{name: "[]byte size:1048576", inputType: []byte{}, size: 1048576, expectTypeStr: "BYTES(1048576)", expectIsNull: false},
		{name: "[]uint8", inputType: []uint8{}, size: 0, expectTypeStr: "BYTES(MAX)", expectIsNull: false},
		{name: "[]uint8 size:1", inputType: []uint8{}, size: 1, expectTypeStr: "BYTES(1)", expectIsNull: false},
		{name: "[]uint8 size:1048576", inputType: []uint8{}, size: 1048576, expectTypeStr: "BYTES(1048576)", expectIsNull: false},
		{name: "time.Time", inputType: time.Time{}, size: 0, expectTypeStr: "TIMESTAMP", expectIsNull: false},
		{name: "spanner.NullTime", inputType: spanner.NullTime{}, size: 0, expectTypeStr: "TIMESTAMP", expectIsNull: true},
		{name: "json.RawMessage", inputType: json.RawMessage{}, size: 0, expectTypeStr: "BYTES(MAX)", expectIsNull: false},
		{name: "string size:0", inputType: "", size: 0, expectTypeStr: "STRING(MAX)", expectIsNull: false},
		{name: "string size:1", inputType: "", size: 1, expectTypeStr: "STRING(1)", expectIsNull: false},
		{name: "string size:2621440", inputType: "", size: 2621440, expectTypeStr: "STRING(2621440)", expectIsNull: false},
		{name: "spanner.NullString size:0", inputType: spanner.NullString{}, size: 0, expectTypeStr: "STRING(MAX)", expectIsNull: true},
		{name: "spanner.NullString size:1", inputType: spanner.NullString{}, size: 1, expectTypeStr: "STRING(1)", expectIsNull: true},
		{name: "spanner.NullString size:2621440", inputType: spanner.NullString{}, size: 2621440, expectTypeStr: "STRING(2621440)", expectIsNull: true},
		{name: "[]bool", inputType: []bool{}, size: 0, expectTypeStr: "ARRAY<BOOL>", expectIsNull: false},
		{name: "[]int8", inputType: []int8{}, size: 0, expectTypeStr: "ARRAY<INT64>", expectIsNull: false},
		{name: "[]int16", inputType: []int16{}, size: 0, expectTypeStr: "ARRAY<INT64>", expectIsNull: false},
		{name: "[]int32", inputType: []int32{}, size: 0, expectTypeStr: "ARRAY<INT64>", expectIsNull: false},
		{name: "[]int", inputType: []int{}, size: 0, expectTypeStr: "ARRAY<INT64>", expectIsNull: false},
		{name: "[]int64", inputType: []int64{}, size: 0, expectTypeStr: "ARRAY<INT64>", expectIsNull: false},
		{name: "[]spanner.NullInt64", inputType: []spanner.NullInt64{}, size: 0, expectTypeStr: "ARRAY<INT64>", expectIsNull: true},
		{name: "[]*spanner.NullInt64", inputType: []*spanner.NullInt64{}, size: 0, expectTypeStr: "ARRAY<INT64>", expectIsNull: true},
		{name: "[]uint16", inputType: []uint16{}, size: 0, expectTypeStr: "ARRAY<INT64>", expectIsNull: false},
		{name: "[]uint32", inputType: []uint32{}, size: 0, expectTypeStr: "ARRAY<INT64>", expectIsNull: false},
		{name: "[]uint", inputType: []uint{}, size: 0, expectTypeStr: "ARRAY<INT64>", expectIsNull: false},
		{name: "[]uint64", inputType: []uint64{}, size: 0, expectTypeStr: "ARRAY<INT64>", expectIsNull: false},
		{name: "[]float32", inputType: []float32{}, size: 0, expectTypeStr: "ARRAY<FLOAT64>", expectIsNull: false},
		{name: "[]float64", inputType: []float64{}, size: 0, expectTypeStr: "ARRAY<FLOAT64>", expectIsNull: false},
		{name: "[]spanner.NullFloat64", inputType: []spanner.NullFloat64{}, size: 0, expectTypeStr: "ARRAY<FLOAT64>", expectIsNull: true},
		{name: "[][]byte", inputType: [][]byte{}, size: 0, expectTypeStr: "ARRAY<BYTES(MAX)>", expectIsNull: false},
		{name: "[][]uint8", inputType: [][]uint8{}, size: 0, expectTypeStr: "ARRAY<BYTES(MAX)>", expectIsNull: false},
		{name: "[]time.Time", inputType: []time.Time{}, size: 0, expectTypeStr: "ARRAY<TIMESTAMP>", expectIsNull: false},
		{name: "[]spanner.NullTime", inputType: []spanner.NullTime{}, size: 0, expectTypeStr: "ARRAY<TIMESTAMP>", expectIsNull: true},
		{name: "[]json.RawMessage", inputType: []json.RawMessage{}, size: 0, expectTypeStr: "ARRAY<BYTES(MAX)>", expectIsNull: false},
		{name: "[]string size:0", inputType: []string{}, size: 0, expectTypeStr: "ARRAY<STRING(MAX)>", expectIsNull: false},
		{name: "[]string size:1", inputType: []string{}, size: 1, expectTypeStr: "ARRAY<STRING(1)>", expectIsNull: false},
		{name: "[]string size:2621440", inputType: []string{}, size: 2621440, expectTypeStr: "ARRAY<STRING(2621440)>", expectIsNull: false},
		{name: "[]spanner.NullString size:0", inputType: []spanner.NullString{}, size: 0, expectTypeStr: "ARRAY<STRING(MAX)>", expectIsNull: true},
		{name: "[]spanner.NullString size:1", inputType: []spanner.NullString{}, size: 1, expectTypeStr: "ARRAY<STRING(1)>", expectIsNull: true},
		{name: "[]spanner.NullString size:2621440", inputType: []spanner.NullString{}, size: 2621440, expectTypeStr: "ARRAY<STRING(2621440)>", expectIsNull: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualTypeStr, actualIsNull := parseTypeToString(reflect.TypeOf(tt.inputType), tt.size)
			if diff := cmp.Diff(tt.expectTypeStr, actualTypeStr); diff != "" {
				t.Errorf("TypeString Diff:\n%s", diff)
			}
			if diff := cmp.Diff(tt.expectIsNull, actualIsNull); diff != "" {
				t.Errorf("IsNull Diff:\n%s", diff)
			}
		})
	}
}

func TestColumn_ToSQL(t *testing.T) {
	type fields struct {
		name        string
		isNull      bool
		size        int
		reflectType reflect.Type
	}
	tests := []struct {
		name   string
		fields fields
		expect string
	}{
		{
			name: "int, not null, no size",
			fields: fields{
				name:        "ID",
				isNull:      false,
				size:        0,
				reflectType: reflect.TypeOf(int64(0)),
			},
			expect: "`ID` INT64 NOT NULL",
		},
		{
			name: "int, nullable, no size",
			fields: fields{
				name:        "ID",
				isNull:      true,
				size:        0,
				reflectType: reflect.TypeOf(int64(0)),
			},
			expect: "`ID` INT64",
		},
		{
			name: "string, not null, size=20",
			fields: fields{
				name:        "Description",
				isNull:      false,
				size:        20,
				reflectType: reflect.TypeOf(""),
			},
			expect: "`Description` STRING(20) NOT NULL",
		},
		{
			name: "string, nullable, size=1024",
			fields: fields{
				name:        "Description",
				isNull:      true,
				size:        1024,
				reflectType: reflect.TypeOf(""),
			},
			expect: "`Description` STRING(1024)",
		},
		{
			name: "byte, not null, size=-1",
			fields: fields{
				name:        "Description",
				isNull:      true,
				size:        -1,
				reflectType: reflect.TypeOf([]byte{}),
			},
			expect: "`Description` BYTES(MAX)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Column{
				name:        tt.fields.name,
				isNull:      tt.fields.isNull,
				size:        tt.fields.size,
				reflectType: tt.fields.reflectType,
			}
			if diff := cmp.Diff(tt.expect, c.ToSQL()); diff != "" {
				t.Errorf("Column.ToSQL() Diff:\n%s", diff)
			}
		})
	}
}
