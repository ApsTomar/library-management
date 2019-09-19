package password_hash

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(p string) (string, error) {
	hash, err := hash(p)
	if err != nil {
		return "", errors.Wrap(err, "error in creating hash of password")
	}
	return hash, nil
}

func hash(password string) (hash string, err error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("error in generating hash of password")
	}
	return string(b), nil
}

func validPassword(password, hash string) bool {
	return nil == bcrypt.CompareHashAndPassword(
		[]byte(hash), []byte(password))
}
