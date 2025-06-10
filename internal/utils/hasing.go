package utils

import "golang.org/x/crypto/bcrypt"

func HashingPassword(password string) string {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	passwordString := string(passwordHash)
	return passwordString
}

func ComparePassword(hashedPwd, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password))
	return err == nil
}
