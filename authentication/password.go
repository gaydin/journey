package authentication

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/gaydin/journey/store"
)

func LoginIsCorrect(ctx context.Context, db store.Database, name string, password string) bool {
	hashedPassword, err := db.RetrieveHashedPasswordForUser(ctx, name)
	if len(hashedPassword) == 0 || err != nil { // len(hashedPassword) == 0 probably not needed.
		// User name likely doesn't exist
		return false
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return false
	}
	return true
}

func EncryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
