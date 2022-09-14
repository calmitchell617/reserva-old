package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"
	"unicode/utf8"

	"github.com/calmitchell617/reserva/internal/validator"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

var AnonymousBank = &Bank{}

type Bank struct {
	Id             int64     `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Password       password  `json:"-"`
	Activated      bool      `json:"activated"`
	Version        int       `json:"-"`
	ExternalId     string    `json:"external_id"`
	BalanceInCents int64     `json:"balance_in_cents"`
	Frozen         bool      `json:"frozen"`
}

func (u *Bank) IsAnonymous() bool {
	return u == AnonymousBank
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(utf8.RuneCountInString(password) >= 8, "password", "must be at least 8 characters long")
	v.Check(utf8.RuneCountInString(password) <= 72, "password", "must not be more than 72 characters long")
}

func ValidateBank(v *validator.Validator, bank *Bank) {
	v.Check(bank.Name != "", "name", "must be provided")
	v.Check(utf8.RuneCountInString(bank.Name) <= 500, "name", "must not be more than 500 characters long")
	v.Check(bank.ExternalId != "", "external_id", "must be provided")

	ValidateEmail(v, bank.Email)

	if bank.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *bank.Password.plaintext)
	}

	if bank.Password.hash == nil {
		panic("missing password hash for bank")
	}
}

type BankModel struct {
	DB *sql.DB
}

func (m BankModel) GetByEmail(email string) (*Bank, error) {
	query := `
        SELECT
					id,
					created_at,
					name,
					email,
					password_hash,
					activated,
					version,
					external_id,
					balance_in_cents,
					frozen
        FROM banks
        WHERE email = $1`

	var bank Bank

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&bank.Id,
		&bank.CreatedAt,
		&bank.Name,
		&bank.Email,
		&bank.Password.hash,
		&bank.Activated,
		&bank.Version,
		&bank.ExternalId,
		&bank.BalanceInCents,
		&bank.Frozen,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &bank, nil
}

func (m BankModel) Update(bank *Bank) error {
	query := `
        UPDATE banks 
        SET
					name = $1,
					email = $2,
					password_hash = $3,
					activated = $4,
					version = version + 1,
					external_id = $5,
					balance_in_cents = $6,
					frozen = $7
        WHERE id = $8 AND version = $9
        RETURNING version`

	args := []interface{}{
		bank.Name,
		bank.Email,
		bank.Password.hash,
		bank.Activated,
		bank.ExternalId,
		bank.BalanceInCents,
		bank.Frozen,
		bank.Id,
		bank.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&bank.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "banks_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m BankModel) GetForToken(tokenScope, tokenPlaintext string) (*Bank, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
        SELECT 
					banks.id,
					banks.created_at,
					banks.name,
					banks.email,
					banks.password_hash,
					banks.activated,
					banks.version,
					banks.external_id,
					banks.balance_in_cents,
					banks.frozen
        FROM banks
        INNER JOIN tokens
        ON banks.id = tokens.bank_id
        WHERE tokens.hash = $1
        AND tokens.scope = $2 
        AND tokens.expiry > $3`

	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var bank Bank

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&bank.Id,
		&bank.CreatedAt,
		&bank.Name,
		&bank.Email,
		&bank.Password.hash,
		&bank.Activated,
		&bank.Version,
		&bank.ExternalId,
		&bank.BalanceInCents,
		&bank.Frozen,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &bank, nil
}
