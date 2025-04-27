package services_testing

import (
	"nexus/services/auth"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {

	// Test the GenerateToken method
	token, err := auth.GenerateToken("test_algorithm", "test_user_id", 3600)
	assert.Equal(t, nil, err, "err should be nil")
	assert.NotEqual(t, "", token, "token should not be empty")

	// Test the ValidateToken method
	validatedUserID, err := auth.ValidateToken(token)
	assert.Equal(t, nil, err, "err should be nil")
	assert.Equal(t, "test_user_id", validatedUserID, "user ID should be 'test_user_id'")
}

func TestExpiredJWT(t *testing.T) {

	// Test the GenerateToken method
	token, err := auth.GenerateToken("test_algorithm", "test_user_id", 1)
	assert.Equal(t, nil, err, "err should be nil")
	assert.NotEqual(t, "", token, "token should not be empty")

	// Test the ValidateToken method
	time.Sleep(1 * time.Second)
	_, err = auth.ValidateToken(token)
	assert.Equal(t, "token has invalid claims: token is expired", err.Error(), "err should be nil")

}
