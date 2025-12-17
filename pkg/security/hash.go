package security

import (
	"golang.org/x/crypto/bcrypt"
)


const cost = bcrypt.DefaultCost // 10

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", nil
	}
	return string(bytes), nil

}

func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}