package auth

import (
	"strconv"

	"github.com/labstack/echo/v4"

	"golang-boilerplate/internal/pkg/utils/rbac"
)

func Middleware(enforcer rbac.RolesManager, obj, act string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := GetUser(c)

			// If is admin, allow access
			if user.IsAdmin {
				return next(c)
			}

			if ok, err := enforcer.Enforce(strconv.Itoa(user.RoleID), obj, act); ok {
				if err != nil {
					return echo.NewHTTPError(403, "Forbidden")
				}
				return next(c)
			}
			return echo.NewHTTPError(403, "Forbidden")
		}
	}
}
