package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Module struct {
	jwtSecret       string
	krakenJwtConfig echojwt.Config
}

func (m *Module) Configure(e *echo.Echo, krakenJwtConfig echojwt.Config, secret string) {
	m.jwtSecret = secret
	m.krakenJwtConfig = krakenJwtConfig

	authRoute := e.Group("/auth/kraken")
	authRoute.Use(echojwt.WithConfig(m.krakenJwtConfig))
	authRoute.POST("", m.AuthKrakenToken)
}

// Authenticate token from Kraken, respond with BAD Token if valid.
func (m *Module) AuthKrakenToken(c echo.Context) error {
	user := GetUser(c)
	claims := &JwtCustomClaims{
		ID:       user.ID,
		Username: user.Username,
		RoleID:   user.RoleID,
		IsAdmin:  user.IsAdmin,
	}

	encToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signToken, err := encToken.SignedString([]byte(m.jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"token": signToken})
}
