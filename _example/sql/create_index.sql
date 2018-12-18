CREATE INDEX `EntryByTitle` ON `Entry` (`Title`)

CREATE NULL_FILTERED INDEX `PlayerCommentByPlayerIDCommentNullFiltered` ON `PlayerComment` (`PlayerID`, `Comment`)

CREATE UNIQUE INDEX `BookmarkByUserIDEntryID` ON `BookMark` (`UserID`, `EntryID` DESC)

CREATE UNIQUE INDEX `BalanceByUserIDCurrencyID` ON `Balance` (`UserID`, `CurrencyID`)