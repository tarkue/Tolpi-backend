package middlewares

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tarkue/tolpi-backend/config"
	"github.com/tarkue/tolpi-backend/internal/app/service"
	usercontext "github.com/tarkue/tolpi-backend/internal/app/userContext"
)

type Middlewares struct {
	s *service.Service
}

func New(s *service.Service) *Middlewares {
	return &Middlewares{
		s: s,
	}
}

func (mw *Middlewares) Authorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		clientUrl := c.Request().Header.Get(AuthHeader)
		if clientUrl == "" {
			c.String(http.StatusOK, AuthError)
			return nil
		}

		if err := mw.s.VerifyLaunchParams(clientUrl, config.SecretKey); err == nil {
			ctx := context.WithValue(c.Request().Context(),
				usercontext.UserCtxKey, &usercontext.UserContext{
					ID: mw.s.GetUserId(c),
				})
			c.SetRequest(c.Request().WithContext(ctx))
			next(c)
			return nil

		} else {
			c.String(http.StatusOK, AuthError)
			return nil
		}

	}
}
