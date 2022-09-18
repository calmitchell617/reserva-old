package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Tokens   TokenModel
	Banks    BankModel
	Accounts AccountModel
	Cards    CardModel
}

func NewModels(writeDb *sql.DB, readDb *sql.DB) Models {
	return Models{
		Tokens:   TokenModel{WriteDb: writeDb, ReadDb: readDb},
		Banks:    BankModel{WriteDb: writeDb, ReadDb: readDb},
		Accounts: AccountModel{WriteDb: writeDb, ReadDb: readDb},
		Cards:    CardModel{WriteDb: writeDb, ReadDb: readDb},
	}
}
