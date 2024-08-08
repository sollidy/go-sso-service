package jwt_test

import (
	"testing"
	"time"

	"sso-service/internal/domain/models"
	"sso-service/internal/lib/jwt"

	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	user := models.User{
		ID:       1,
		Email:    "test@example.com",
		PassHash: []byte("password"),
	}
	app := models.App{
		ID:     1,
		Secret: "secret",
	}
	duration := time.Hour

	token, err := jwt.NewToken(user, app, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// claims, err := jwt.ParseToken(token)
	// assert.NoError(t, err)
	// assert.Equal(t, user.ID, claims["uid"])
	// assert.Equal(t, user.Email, claims["email"])
	// assert.Equal(t, app.ID, claims["app_id"])
	// assert.WithinDuration(t, time.Now().Add(duration), time.Unix(int64(claims["exp"].(float64)), 0), time.Second)
}

// func TestNewToken_Error(t *testing.T) {
// 	user := models.User{
// 		ID:       1,
// 		Email:    "test@example.com",
// 		PassHash: []byte("password"),
// 	}
// 	app := models.App{
// 		ID:     1,
// 		Secret: "",
// 	}
// 	duration := time.Hour

// 	_, err := jwt.NewToken(user, app, duration)
// 	assert.Error(t, err)
// }
