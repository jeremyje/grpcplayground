package auth

import "fmt"

const (
	secretTokenValue = "the-secret-token"
)

// VerifyToken verifies that the token is authentic.
func VerifyToken(token string) error {
	if token == secretTokenValue {
		return nil
	}
	return fmt.Errorf("token [%s] does not match the secret token value", token)
}
