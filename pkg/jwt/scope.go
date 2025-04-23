package jwt

import (
	"encoding/base64"
	"encoding/json"

	"github.com/nguyenthanhtrung001/reup/internal/models"
)

// NewScope creates a new scope.
func NewScope(payload Payload) models.Scope {
	return models.Scope{
		UserID:    payload.UserID,
		GroupID:   payload.GroupID,
		GroupRole: payload.GroupRole,
	}
}

func CreateScopeHeader(scope models.Scope) (string, error) {
	jsonData, err := json.Marshal(scope)
	if err != nil {
		return "", err
	}

	// Encode the JSON data as Base64
	base64Data := base64.StdEncoding.EncodeToString(jsonData)

	return base64Data, nil
}

func ParseScopeHeader(scopeHeader string) (models.Scope, error) {
	// Decode the Base64 data
	jsonData, err := base64.StdEncoding.DecodeString(scopeHeader)
	if err != nil {
		return models.Scope{}, err
	}

	// Unmarshal the JSON data
	var scope models.Scope
	err = json.Unmarshal(jsonData, &scope)
	if err != nil {
		return models.Scope{}, err
	}

	return scope, nil
}
