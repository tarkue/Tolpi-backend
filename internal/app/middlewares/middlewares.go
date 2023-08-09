package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tarkue/tolpi-backend/config"
	"github.com/tarkue/tolpi-backend/internal/app/service"
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
		next(c)
		return nil

		clientUrl := c.Request().Header.Get("Authorization")
		if clientUrl == "" {
			c.String(http.StatusOK, "error auth")
			return nil
		}

		if err := mw.s.VerifyLaunchParams(clientUrl, config.SecretKey); err == nil {
			next(c)
			return nil

		} else {
			c.String(http.StatusOK, "error auth")
			return nil
		}

	}
}
