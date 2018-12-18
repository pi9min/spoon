package example

import (
	"time"

	"cloud.google.com/go/civil"
	"cloud.google.com/go/spanner"
	"github.com/pi9min/spoon"
)

var (
	_ spoon.EntityBehavior = (*User)(nil)
	_ spoon.EntityBehavior = Entry{}
	_ spoon.EntityBehavior = PlayerComment{}
	_ spoon.EntityBehavior = Bookmark{}
	_ spoon.EntityBehavior = Balance{}
)

type User struct {
	ID         uint64
	Name       string
	Token      string
	BornedDate spanner.NullDate `db:"nullable"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (u *User) TableName() string {
	return "User"
}

func (u *User) PrimaryKey() *spoon.PrimaryKey {
	return spoon.AddPrimaryKey(spoon.KeyPart{ColumnName: "ID"})
}

func (u *User) Indexes() spoon.Indexes {
	return spoon.Indexes{}
}

type Entry struct {
	ID        int32
	Title     string
	Public    bool
	Content   *string `db:"size=1048576"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e Entry) TableName() string {
	return "Entry"
}

func (e Entry) PrimaryKey() *spoon.PrimaryKey {
	return spoon.AddPrimaryKeyWithInterleave("User", spoon.KeyPart{ColumnName: "ID"}, spoon.KeyPart{ColumnName: "CreatedAt", IsOrderDesc: true})
}

func (e Entry) Indexes() spoon.Indexes {
	return spoon.Indexes{
		spoon.AddIndex("EntryByTitle", "Entry", false, spoon.KeyPart{ColumnName: "Title"}),
	}
}

type PlayerComment struct {
	ID        int32              `json:"id"`
	PlayerID  int32              `json:"player_id"`
	EntryID   int32              `json:"entry_id"`
	Comment   spanner.NullString `db:"size:99, nullable" json:"comment"`
	CreatedAt time.Time          `json:"created_at"`
	updatedAt time.Time
}

func (pc PlayerComment) TableName() string {
	return "PlayerComment"
}

func (pc PlayerComment) PrimaryKey() *spoon.PrimaryKey {
	return spoon.AddPrimaryKey(spoon.KeyPart{ColumnName: "ID"})
}

func (pc PlayerComment) Indexes() spoon.Indexes {
	return spoon.Indexes{
		spoon.AddIndex("PlayerCommentByPlayerIDCommentNullFiltered", "PlayerComment", true, spoon.KeyPart{ColumnName: "PlayerID"}, spoon.KeyPart{ColumnName: "Comment"}),
	}
}

type Bookmark struct {
	ID        string
	UserID    int32
	EntryID   int32
	Ignore    string `db:"-"`
	Comments  []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b Bookmark) TableName() string {
	return "Bookmark"
}

func (b Bookmark) PrimaryKey() *spoon.PrimaryKey {
	return spoon.AddPrimaryKey(spoon.KeyPart{ColumnName: "ID"})
}

func (b Bookmark) Indexes() spoon.Indexes {
	return spoon.Indexes{
		spoon.AddUniqueIndex("BookmarkByUserIDEntryID", "BookMark", false, spoon.KeyPart{ColumnName: "UserID"}, spoon.KeyPart{ColumnName: "EntryID", IsOrderDesc: true}),
	}
}

// protobufのenum扱いのパターン
type CurrencyID int32
type Balance struct {
	ID         string
	UserID     string
	CurrencyID CurrencyID
	Amount     float64
}

func (b Balance) TableName() string {
	return "Balance"
}

func (b Balance) PrimaryKey() *spoon.PrimaryKey {
	return spoon.AddPrimaryKey(spoon.KeyPart{ColumnName: "ID"})
}

func (b Balance) Indexes() spoon.Indexes {
	return spoon.Indexes{
		spoon.AddUniqueIndex("BalanceByUserIDCurrencyID", "Balance", false, spoon.KeyPart{ColumnName: "UserID"}, spoon.KeyPart{ColumnName: "CurrencyID"}),
	}
}

type NestChild1 struct {
	NC1ID    string
	NestedAt time.Time
}

type NestChild2 struct {
	NC2ID       string
	Birthdate   civil.Date
	IgnoreField []byte `db:"-"`
	Nested2At   spanner.NullTime
}

type NestParent struct {
	*NestChild1
	*NestChild2
}

func (n NestParent) TableName() string {
	return "NestParent"
}

func (n NestParent) PrimaryKey() *spoon.PrimaryKey {
	return spoon.AddPrimaryKey(spoon.KeyPart{ColumnName: "NC1ID"})
}

func (n NestParent) Indexes() spoon.Indexes {
	return spoon.Indexes{}
}
