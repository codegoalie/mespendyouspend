package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// SpenderToken is used by pop to map your spender_tokens database table to your go code.
type SpenderToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	SpenderID uuid.UUID `json:"spender_id" db:"spender_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Spender Spender `belongs_to:"spenders" json:"spender" db:"-"`
}

// String is not required by pop and may be deleted
func (s SpenderToken) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// SpenderTokens is not required by pop and may be deleted
type SpenderTokens []SpenderToken

// String is not required by pop and may be deleted
func (s SpenderTokens) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *SpenderToken) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.Token, Name: "Token"},
		&validators.TimeIsPresent{Field: s.ExpiresAt, Name: "ExpiresAt"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *SpenderToken) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *SpenderToken) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
