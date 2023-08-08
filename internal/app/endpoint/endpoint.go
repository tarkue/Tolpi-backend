package endpoint

import (
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tarkue/tolpi-backend/config"
	"github.com/tarkue/tolpi-backend/internal/app/graph/model"
)

type DataBase interface {
	FindUserById(userID string) *model.User
	UpdateUserTrackers(userID string, trackers []string)
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
	if !ctx.QueryParams().Has("user_id") {
		ctx.String(http.StatusOK, "not found user or sub_id")
		return nil
	}

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

func (e *Endpoint) Subscribe(ctx echo.Context) error {
	if !ctx.QueryParams().Has("user_id") || !ctx.QueryParams().Has("sub_id") {
		ctx.String(http.StatusOK, "not found user or sub_id")
		return nil
	}

	userId := ctx.QueryParam("user_id")
	subId := ctx.QueryParam("sub_id")
	user := e.db.FindUserById(subId)

	NewSlice := append(user.TrackerList, userId)
	e.db.UpdateUserTrackers(subId, NewSlice)

	ctx.String(http.StatusOK, "true")
	return nil
}

func (e *Endpoint) Unsubscribe(ctx echo.Context) error {
	if !ctx.QueryParams().Has("user_id") || !ctx.QueryParams().Has("unsub_id") {
		ctx.String(http.StatusOK, "not found user or sub_id")
		return nil
	}

	userId := ctx.QueryParam("user_id")
	unsubId := ctx.QueryParam("unsub_id")
	user := e.db.FindUserById(unsubId)

	NewSlice := RemoveIndex(user.TrackerList, IndexOf(userId, user.TrackerList))
	e.db.UpdateUserTrackers(unsubId, NewSlice)

	ctx.String(http.StatusOK, "true")
	return nil
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}
