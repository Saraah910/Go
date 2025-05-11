package Utils

import "golang.org/x/crypto/bcrypt"

func ConvertToHash(password string) (string, error) {
	byteArray, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(byteArray), err
}

func CheckPassword(userInput string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userInput))
	return err == nil
}
