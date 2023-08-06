package endpoint

import (
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tarkue/tolpi-backend/config"
)

type DataBase interface {
}

type Endpoint struct {
	db DataBase
}

func New(db DataBase) *Endpoint {
	return &Endpoint{
		db: db,
	}
}

func (e *Endpoint) GetStatus(ctx echo.Context) error {
	link := config.VkApiLink + config.VkUsersGetMethod + "?" + `access_token=` + config.VkServiceToken + `&user_ids=` + ctx.QueryParam("user_id") + `&fields=status&v=5.131`

	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Print(link)

	ctx.String(http.StatusOK, string(body))
	return nil
}

func (e *Endpoint) GetCountry(ctx echo.Context) error {

	resp, err := http.Get(config.CountriesApi)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ctx.String(http.StatusOK, string(body))
	return nil
}

func (e *Endpoint) GetTrackers(ctx echo.Context) error {
	return nil
}
