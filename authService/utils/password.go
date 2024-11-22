package utils

import "golang.org/x/crypto/bcrypt"


func HashPasswordWithSalt(password, salt string) (string, error) {
	passwordWithSalt := password + salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordWithSalt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordWithSalt(hashedPassword, password, salt string) bool {
	passwordWithSalt := password + salt
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(passwordWithSalt))
	return err == nil
}