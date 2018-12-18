# spoon

`spoon` is a library to generate [Google Cloud Spanner](https://cloud.google.com/spanner/) table schema.

`spoon` parses the structure and generates the table schema converted to the appropriate type.

## How to use

### 1. Define the table structure

Please implement the structure to satisfy the following `spoon.EntityBehavior` interface.

```go
type EntityBehavior interface {
	TableName() string
	PrimaryKey() spoon.PrimaryKey
	Indexes() spoon.Indexes
}
```

`_example/readme/user/user.go` See the source code below.

```go
package user

import (
	"time"

	"github.com/pi9min/spoon"
)

// validation
var _ spoon.EntityBehavior = (*User)(nil)

type Sex int8

const (
	Male Sex = iota
	Famale
)

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Sex       Sex
	Age       int64
	FriendIDs []int64
	CreatedAt time.Time
}

func (u *User) TableName() string {
	return "User"
}

func (u *User) PrimaryKey() *spoon.PrimaryKey {
	return spoon.AddPrimaryKey(spoon.KeyPart{ColumnName: "ID"})
}

func (u *User) Indexes() spoon.Indexes {
	return spoon.Indexes{
		spoon.AddIndex(
			"UserByLastFirstName",
			"User",
			false,
			spoon.KeyPart{ColumnName: "LastName"},
			spoon.KeyPart{ColumnName: "FirstName"},
		),
	}
}
```


### 2. Generate spoon client, load structure and output table schema

`_example/readme/main.go` See the source code below.

```go
package main

import (
	"fmt"
	"os"

	"github.com/pi9min/spoon"
	"github.com/pi9min/spoon/_example/readme/user"
)

func main() {
	cli, err := spoon.New()
	if err != nil {
		panic(err)
	}

	schema, err := cli.GenerateCreateTable(&user.User{})
	if err != nil {
		panic(err)
	}

	// output to stdout
	fmt.Fprint(os.Stdout, schema)
}
```

The execution result is as follows.

```bash
$ go run ./_example/readme/main.go
CREATE TABLE `User` (
    `ID` INT64 NOT NULL,
    `FirstName` STRING(MAX) NOT NULL,
    `LastName` STRING(MAX) NOT NULL,
    `Sex` INT64 NOT NULL,
    `Age` INT64 NOT NULL,
    `FriendIDs` ARRAY<INT64> NOT NULL,
    `CreatedAt` TIMESTAMP NOT NULL,
) PRIMARY KEY (`ID`)
```

## A type correspondence table between spanner and golang

|        Golang Type       | Spanner Type |
| :----------------------: | :----------: |
|           int8           |    `INT64`     |
|           int16          |    `INT64`     |
|           int32          |    `INT64`     |
| int64, spanner.NullInt64 |    `INT64`     |
|           uint8          |    `INT64`     |
|          uint16          |    `INT64`     |
|          uint32          |    `INT64`     |
|          uint64          |    `INT64`     |
|         float32          |   `FLOAT64`    |
|      []byte,[]uint8      |    `BYTES(n or MAX) (1 <= n <= 10485760)`     |
| float64, spanner.NullFloat64 |   `FLOAT64`    |
|  string, spanner.NullString  |  `STRING(n or MAX) (1 <= n <= 2621440)`   |
|    bool, spanner.NullBool    |    `BOOL`      |
|    civil.Date, spanner.NullDate    |    `DATE`      |
| time.Time, spanner.NullTime        |  `TIMESTAMP`   |
|     json.RawMessage      |   `BYTES(n)`   |
|     Primitive type slices  |  `ARRAY<TYPE>` |

## Structure tag prefix

Default tag prefix is `db`.

You can set your favorite prefix by specifying TagPrefix option.

e.g. Set `spanner` to prefix.

```go
	cli, err := spoon.New(TagPrefix("spanner"))
	if err != nil {
		panic(err)
	}
```

Values that can be specified with the tag are as follows.

|   Tag Value   |                       VALUE                       |
| :-----------: | :-----------------------------------------------: |
|   `nullable`    |   Remove the `NOT NULL` constraint (Allow NULL)   |
|   `size=<n>`  |   When it is strings or bytes, set the length     |
|      `-`        |                   Ignore fields                   |

It's used as follows.

```go
type User struct {
	ID            string `db:"size=64"`
	Name          string
	TemporaryMemo string `db:"-"`
	CreatedAt     time.Time
	DeletedAt     spanner.NullTime `db:"nullable"`
}
```

## How to set the PrimaryKey

Uses `spoon.AddPrimaryKey()` method.

If you also want to set interleaving, use `spoon.AddPrimaryKeyWithInterleave()` method.

For example, 
```go
// Set ID as PrimaryKey
func (u *User) PrimaryKey() *spoon.PrimaryKey {
	return spoon.AddPrimaryKey(spoon.KeyPart{ColumnName:"ID"})
}

---> PRIMARY KEY (`ID`)


// Set ID and CreatedAt DESC to PrimaryKey
func (u *User) PrimaryKey() *spoon.PrimaryKey {
	return spoon.AddPrimaryKey(
		spoon.KeyPart{ColumnName: "ID"},
		spoon.KeyPart{ColumnName: "CreatedAt", IsOrderDesc: true},
	)
}

--> PRIMARY KEY (`ID`, `CreatedAt` DESC)


// Interleave to Company table and set ID to PrimaryKey
func (u *User) PrimaryKey() *spoon.PrimaryKey {
	return spoon.AddPrimaryKeyWithInterleave("Company", spoon.KeyPart{ColumnName: "ID"})
}

--> PRIMARY KEY (`ID`), INTERLEAVE IN PARENT `Company`
```

## How to set the Index

Uses `spoon.AddIndex()` method.

If you also want to add a unique constraint, use `spoon.AddUniqueIndex()` method.

For example,
```go
// Add index to Lastname, FirstName
func (u *User) Indexes() spoon.Indexes {
	return spoon.Indexes{
		spoon.AddIndex(
			"UserByLastFirstName",
			"User",
			false,
			spoon.KeyPart{ColumnName: "LastName"},
			spoon.KeyPart{ColumnName: "FirstName"},
		),
	}
}

--> CREATE INDEX `UserByLastFirstName` ON `User` (`LastName`, `FirstName`)

// Add a unique index to Lastname, FirstName
func (u *User) Indexes() spoon.Indexes {
	return spoon.Indexes{
		spoon.AddUniqueIndex(
			"UserByUniqueLastFirstName",
			"User",
			false,
			spoon.KeyPart{ColumnName: "LastName"},
			spoon.KeyPart{ColumnName: "FirstName"},
		),
	}
}

--> CREATE UNIQUE INDEX `UserByUniqueLastFirstName` ON `User` (`LastName`, `FirstName`)

// Add an index to CreatedAt DESC without including NULL
func (u *User) Indexes() spoon.Indexes {
	return spoon.Indexes{
		spoon.AddIndex(
			"UserByNullFilteredCreatedAtDesc",
			"User",
			true,
			spoon.KeyPart{ColumnName: "CreatedAt", IsOrderDesc: true},
		),
	}
}

--> CREATE NULL_FILTERED INDEX `UserByNullFilteredCreatedAtDesc` ON `User` (`CreatedAt` DESC)
```

## License

See [LICENSE.md](/LICENSE.md)
