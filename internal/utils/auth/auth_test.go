package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/stretchr/testify/assert"
)

func TestInitJwtAuth(t *testing.T) {
	type args struct {
		secret string
	}
	tests := []struct {
		name string
		args args
		want echojwt.Config
	}{
		{
			name: "Test Init",
			args: args{
				secret: "abcd",
			},
			want: echojwt.Config{
				NewClaimsFunc: func(_ echo.Context) jwt.Claims {
					return new(JwtCustomClaims)
				},
				SigningKey: []byte("abcd"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InitJwtAuth(tt.args.secret)
			assert.Equal(t, got.SigningKey, tt.want.SigningKey)
		})
	}
}

func TestGetUser(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()
	secret := "secret"

	// Create a sample JWT token with claims
	claims := &JwtCustomClaims{
		ID:       123,
		Username: "testuser",
		RoleID:   456,
		IsAdmin:  true,
	}
	signToken, encToken := generateDummyToken(claims, secret)

	// Create a new HTTP request with the token in the context
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+signToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", encToken)

	// Call the GetUser function
	user := GetUser(c)
	// Assert that the user information is correctly extracted
	assert.Equal(t, 123, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, 456, user.RoleID)
	assert.True(t, user.IsAdmin)

	// Set with wrong secret token
	secret = "wrongsecret"
	signToken, _ = generateDummyToken(claims, secret)
	req.Header.Set("Authorization", "Bearer "+signToken)
	cx := e.NewContext(req, rec)
	user = GetUser(cx)
	// Assert that the user information is correctly extracted
	assert.NotEqual(t, 123, user.ID)
	assert.Equal(t, 0, user.ID)
	assert.NotEqual(t, "testuser", user.Username)
	assert.NotEqual(t, 456, user.RoleID)
	assert.False(t, user.IsAdmin)
}

// Helper function to create a sample JWT token for testing.
func generateDummyToken(claims *JwtCustomClaims, secret string) (string, *jwt.Token) {
	encToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signToken, _ := encToken.SignedString([]byte(secret))

	return signToken, encToken
}
