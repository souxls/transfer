package password

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(pwd string) string {
	encryptPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(encryptPwd)
}

func ValidatePassword(encryptPassword string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(encryptPassword), []byte(password)); err != nil {
		return false
	}
	return true
}
