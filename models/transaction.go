package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Transaction is used by pop to map your transactions database table to your go code.
type Transaction struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Description string    `json:"description" db:"description"`
	Amount      int       `json:"amount" db:"amount"`
	SpenderID   uuid.UUID `json:"spender_id" db:"spender_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	Spender *Spender `belongs_to:"spenders" json:"spender" db:"-"`
}

// String is not required by pop and may be deleted
func (t Transaction) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Transactions is not required by pop and may be deleted
type Transactions []Transaction

// String is not required by pop and may be deleted
func (t Transactions) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Transaction) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.Description, Name: "Description"},
		&validators.IntIsPresent{Field: t.Amount, Name: "Amount"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Transaction) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Transaction) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
