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
}

func NewModels(db *sql.DB) Models {
	return Models{
		Tokens:   TokenModel{DB: db},
		Banks:    BankModel{DB: db},
		Accounts: AccountModel{DB: db},
	}
}
