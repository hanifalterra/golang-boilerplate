package auth

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	RoleID   int    `json:"role_id"`
	IsAdmin  bool   `json:"is_admin"`
}

type JwtCustomClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	RoleID   int    `json:"role_id"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// Valid implements jwt.Claims.
func (*JwtCustomClaims) Valid() error {
	return nil
}

func InitJwtAuth(secret string) echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(_ echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey: []byte(secret),
	}
}

func GetUser(c echo.Context) User {
	userInfo := User{}
	if user, ok := c.Get("user").(*jwt.Token); ok {
		claims, _ := user.Claims.(*JwtCustomClaims)
		userInfo.ID = claims.ID
		userInfo.Username = claims.Username
		userInfo.RoleID = claims.RoleID
		userInfo.IsAdmin = claims.IsAdmin
		return userInfo
	}
	return userInfo
}
