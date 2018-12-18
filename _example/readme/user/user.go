package user

import (
	"time"

	"github.com/pi9min/spoon"
)

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
