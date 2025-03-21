package jwt

import (
	"encoding/base64"
	"encoding/json"

	"book-store/internal/models"
	pkgCrt "book-store/pkg/encrypter"
)

func CreateApiKey(scope models.Scope, ecnrypter pkgCrt.Encrypter) (string, error) {
	jsonData, err := json.Marshal(scope)
	if err != nil {
		return "", err
	}
	restKey, err := ecnrypter.Encrypt(string(jsonData))
	if err != nil {
		return "", err
	}
	// Encode the JSON data as Base64
	base64Data := base64.StdEncoding.EncodeToString([]byte(restKey))

	return base64Data, nil
}

func ParseApiKey(restKey string, ecnrypter pkgCrt.Encrypter) (models.Scope, error) {
	// Decode the Base64 data
	jsonData, err := base64.StdEncoding.DecodeString(restKey)
	if err != nil {
		return models.Scope{}, err
	}
	restKey, err = ecnrypter.Decrypt(string(jsonData))
	if err != nil {
		return models.Scope{}, err
	}
	// Unmarshal the JSON data
	var scope models.Scope
	err = json.Unmarshal([]byte(restKey), &scope)
	if err != nil {
		return models.Scope{}, err
	}

	return scope, nil
}
