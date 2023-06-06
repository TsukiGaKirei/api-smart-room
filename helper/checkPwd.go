package helper

import (
	"golang.org/x/crypto/bcrypt"
)

func IsPwdFalse(hashedPwd, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd))
	return err != nil
}
