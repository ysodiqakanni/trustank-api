package utility

import "golang.org/x/crypto/bcrypt"

func hashPasswordString(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}
