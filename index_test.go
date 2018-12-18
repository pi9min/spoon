package spoon_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pi9min/spoon"
)

func TestAddIndex(t *testing.T) {
	tests := []struct {
		name   string
		index  *spoon.Index
		expect string
	}{
		{
			name:   "Player.PlayerID, nullFiltered=false",
			index:  spoon.AddIndex("PlayerByPlayerID", "Player", false, spoon.KeyPart{ColumnName: "PlayerID"}),
			expect: "CREATE INDEX `PlayerByPlayerID` ON `Player` (`PlayerID`)",
		},
		{
			name:   "Player.PlayerID, nullFiltered=false",
			index:  spoon.AddIndex("PlayerByPlayerIDEntryID", "Player", true, spoon.KeyPart{ColumnName: "PlayerID", IsOrderDesc: true}, spoon.KeyPart{ColumnName: "EntryID"}),
			expect: "CREATE NULL_FILTERED INDEX `PlayerByPlayerIDEntryID` ON `Player` (`PlayerID` DESC, `EntryID`)",
		},
		{
			name:   "Player.PlayerID, nullFiltered=false",
			index:  spoon.AddUniqueIndex("PlayerByUniquePlayerID", "Player", false, spoon.KeyPart{ColumnName: "PlayerID"}),
			expect: "CREATE UNIQUE INDEX `PlayerByUniquePlayerID` ON `Player` (`PlayerID`)",
		},
		{
			name:   "Player.PlayerID, nullFiltered=false",
			index:  spoon.AddUniqueIndex("PlayerByUniquePlayerIDEntryID", "Player", true, spoon.KeyPart{ColumnName: "PlayerID", IsOrderDesc: true}, spoon.KeyPart{ColumnName: "EntryID"}),
			expect: "CREATE UNIQUE NULL_FILTERED INDEX `PlayerByUniquePlayerIDEntryID` ON `Player` (`PlayerID` DESC, `EntryID`)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.index.CreateIndexSchema()
			if diff := cmp.Diff(tt.expect, actual); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}
