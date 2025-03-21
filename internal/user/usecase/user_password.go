package usecase

import (
	"crypto/sha256"
	"encoding/hex"
)

func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))
	return hashedPassword
}

func checkPassword(password, hashedPassword string) bool {
	return hashPassword(password) == hashedPassword
}
