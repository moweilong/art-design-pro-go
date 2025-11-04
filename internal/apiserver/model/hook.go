package model

import (
	"github.com/google/uuid"
	"github.com/moweilong/art-design-pro-go/internal/pkg/known"
	"github.com/onexstack/onexstack/pkg/authn"
	"github.com/onexstack/onexstack/pkg/rid"
	"github.com/onexstack/onexstack/pkg/store/registry"
	"gorm.io/gorm"
)

// BeforeCreate runs before creating a SecretM database record and initializes various fields.
func (m *SecretM) BeforeCreate(tx *gorm.DB) (err error) {
	// Supports custom SecretKey, but users need to ensure the uniqueness of the SecretKey themselves.
	// onex-cacheserver will use this feature to set secret.
	if m.SecretID == "" {
		// Generate a new UUID for SecretKey.
		m.SecretID = uuid.New().String()
	}

	// Generate a new UUID for SecretID.
	m.SecretKey = uuid.New().String()

	// Set the default status for the secret as normal.
	m.Status = known.SecretStatusNormal

	return nil
}

// BeforeCreate encrypts the plaintext password before creating a database record.
func (m *UserM) BeforeCreate(tx *gorm.DB) error {
	// Encrypt the user password.
	var err error
	m.Password, err = authn.Encrypt(m.Password)
	if err != nil {
		return err
	}

	return nil
}

// AfterCreate generates a userID after creating a database record.
func (m *UserM) AfterCreate(tx *gorm.DB) error {
	m.UserID = rid.NewResourceID("user").New(uint64(m.ID))

	return tx.Save(m).Error
}

func init() {
	registry.Register(&UserM{})
}
