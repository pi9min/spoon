package spoon_test

import (
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/go-cmp/cmp"
	"github.com/pi9min/spoon"
)

var (
	_ spoon.EntityBehavior = Test1{}
	_ spoon.EntityBehavior = (*Test2)(nil)
)

type Test1 struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t1 Test1) TableName() string {
	return "Test1"
}

func (t1 Test1) Indexes() spoon.Indexes {
	return spoon.Indexes{
		spoon.AddIndex("Test1ByCreatedAtDesc", "Test1", false, spoon.KeyPart{ColumnName: "CreatedAt", IsOrderDesc: true}),
	}
}

func (t1 Test1) PrimaryKey() *spoon.PrimaryKey {
	return spoon.AddPrimaryKey(spoon.KeyPart{ColumnName: "ID"})
}

type Test2 struct {
	ID        uint64
	Test1ID   uint64
	Comment   spanner.NullString `db:"size=1024"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t2 *Test2) TableName() string {
	return "Test2"
}

func (t2 *Test2) Indexes() spoon.Indexes {
	return spoon.Indexes{}
}

func (t2 *Test2) PrimaryKey() *spoon.PrimaryKey {
	return spoon.AddPrimaryKeyWithInterleave("Test1", spoon.KeyPart{ColumnName: "ID"}, spoon.KeyPart{ColumnName: "CreatedAt", IsOrderDesc: true})
}

func TestGenerateCreateTable(t *testing.T) {
	tests := []struct {
		name   string
		entity spoon.EntityBehavior
		expect string
	}{
		{
			name:   "1 all not null",
			entity: &Test1{},
			expect: fmt.Sprintf(`CREATE TABLE %s (
    %s INT64 NOT NULL,
    %s STRING(MAX) NOT NULL,
    %s TIMESTAMP NOT NULL,
    %s TIMESTAMP NOT NULL,
) PRIMARY KEY (%s)`,
				spoon.Quote("Test1"),
				spoon.Quote("ID"),
				spoon.Quote("Name"),
				spoon.Quote("CreatedAt"),
				spoon.Quote("UpdatedAt"),
				spoon.Quote("ID"),
			),
		},
		{
			name:   "2 use nullable, size, with interleave",
			entity: &Test2{},
			expect: fmt.Sprintf(`CREATE TABLE %s (
    %s INT64 NOT NULL,
    %s INT64 NOT NULL,
    %s STRING(1024),
    %s TIMESTAMP NOT NULL,
    %s TIMESTAMP NOT NULL,
) PRIMARY KEY (%s, %s), INTERLEAVE IN PARENT %s`,
				spoon.Quote("Test2"),
				spoon.Quote("ID"),
				spoon.Quote("Test1ID"),
				spoon.Quote("Comment"),
				spoon.Quote("CreatedAt"),
				spoon.Quote("UpdatedAt"),
				spoon.Quote("ID"), spoon.Quote("CreatedAt")+" DESC", spoon.Quote("Test1"),
			),
		},
	}

	cli, err := spoon.New()
	if err != nil {
		t.Fatalf("error new Client")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := cli.GenerateCreateTable(tt.entity)
			if err != nil {
				t.Fatalf("error generate create table schema %#v", err)
			}

			if diff := cmp.Diff(tt.expect, actual); diff != "" {
				t.Errorf("GenerateCreateTable Diff:\n%s", diff)
			}
		})
	}
}
