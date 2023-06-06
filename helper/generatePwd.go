package helper

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePwd(origin string) string {
	originToByte := []byte(origin)
	hashedPassword, _ := bcrypt.GenerateFromPassword(originToByte, bcrypt.DefaultCost)

	return string(hashedPassword)
}
