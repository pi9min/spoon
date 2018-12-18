package spoon_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pi9min/spoon"
)

func TestAddPrimaryKey(t *testing.T) {
	tests := []struct {
		name   string
		pk     *spoon.PrimaryKey
		expect string
	}{
		{
			name:   "1 single key part",
			pk:     spoon.AddPrimaryKey(spoon.KeyPart{ColumnName: "ID"}),
			expect: "PRIMARY KEY (`ID`)",
		},
		{
			name:   "2 multiple key part",
			pk:     spoon.AddPrimaryKey(spoon.KeyPart{ColumnName: "ID"}, spoon.KeyPart{ColumnName: "CreatedAt", IsOrderDesc: true}),
			expect: "PRIMARY KEY (`ID`, `CreatedAt` DESC)",
		},
		{
			name:   "3 multiple key part with interleave",
			pk:     spoon.AddPrimaryKeyWithInterleave("Balance", spoon.KeyPart{ColumnName: "CurrencyID"}, spoon.KeyPart{ColumnName: "UserID", IsOrderDesc: true}),
			expect: "PRIMARY KEY (`CurrencyID`, `UserID` DESC), INTERLEAVE IN PARENT `Balance`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.pk.ToSQL()
			if diff := cmp.Diff(tt.expect, actual); diff != "" {
				t.Errorf("ToSQL Diff:\n%s", diff)
			}
		})
	}
}
