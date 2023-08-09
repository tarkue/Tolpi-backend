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
	UpdateUserAvatar(userID string, avatar string)
	UpdateUserFirstName(userID string, firstName string)
	UpdateUserLastName(userID string, lastName string)
	UpdateUserTrackers(userID string, trackers []string)
}

type Service interface {
	RemoveIndex(s []string, index int) []string
	IndexOf(element string, data []string) int
	Contains(a []string, x string) bool
	GetUserId(c echo.Context) string
}

type Endpoint struct {
	db DataBase
	s  Service
}

func New(db DataBase, s Service) *Endpoint {
	return &Endpoint{
		db: db,
		s:  s,
	}
}

func (e *Endpoint) GetStatus(ctx echo.Context) error {
	userId := e.s.GetUserId(ctx)

	link := config.VkApiLink + config.VkUsersGetMethod + "?" + `access_token=` + config.VkServiceToken + `&user_ids=` + userId + `&fields=status&v=5.131`

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

func (e *Endpoint) CheckSubscribe(ctx echo.Context) error {
	if ctx.QueryParams().Has(subIdParam) {
		ctx.String(http.StatusOK, notFoundMessage+subIdParam)
		return nil
	}

	userId := e.s.GetUserId(ctx)
	subId := ctx.QueryParam(subIdParam)
	user := e.db.FindUserById(subId)

	ctx.JSON(http.StatusOK, e.s.Contains(user.TrackerList, userId))

	return nil
}

func (e *Endpoint) Subscribe(ctx echo.Context) error {
	if !ctx.QueryParams().Has(subIdParam) {
		ctx.String(http.StatusOK, notFoundMessage+subIdParam)
		return nil
	}

	userId := e.s.GetUserId(ctx)
	subId := ctx.QueryParam(subIdParam)
	user := e.db.FindUserById(subId)

	NewSlice := append(user.TrackerList, userId)
	e.db.UpdateUserTrackers(subId, NewSlice)

	ctx.JSON(http.StatusOK, true)
	return nil
}

func (e *Endpoint) Unsubscribe(ctx echo.Context) error {
	if !ctx.QueryParams().Has(unSubIdParam) {
		ctx.String(http.StatusOK, notFoundMessage+unSubIdParam)
		return nil
	}

	userId := e.s.GetUserId(ctx)
	unsubId := ctx.QueryParam(unSubIdParam)
	user := e.db.FindUserById(unsubId)

	NewSlice := e.s.RemoveIndex(user.TrackerList, e.s.IndexOf(userId, user.TrackerList))
	e.db.UpdateUserTrackers(unsubId, NewSlice)

	ctx.JSON(http.StatusOK, true)
	return nil
}

func (e *Endpoint) SetAvatar(ctx echo.Context) error {
	if !ctx.QueryParams().Has(avatarParam) {
		ctx.String(http.StatusOK, notFoundMessage+avatarParam)
		return nil
	}
	avatar := ctx.QueryParam(avatarParam)
	userId := e.s.GetUserId(ctx)

	e.db.UpdateUserAvatar(userId, avatar)
	ctx.JSON(http.StatusOK, true)
	return nil
}

func (e *Endpoint) SetFirstName(ctx echo.Context) error {
	if !ctx.QueryParams().Has(firstNameParam) {
		ctx.String(http.StatusOK, notFoundMessage+firstNameParam)
		return nil
	}
	firstName := ctx.QueryParam(firstNameParam)
	userId := e.s.GetUserId(ctx)

	e.db.UpdateUserFirstName(userId, firstName)
	ctx.JSON(http.StatusOK, true)
	return nil
}

func (e *Endpoint) SetLastName(ctx echo.Context) error {
	if !ctx.QueryParams().Has(lastNameParam) {
		ctx.String(http.StatusOK, notFoundMessage+lastNameParam)
		return nil
	}
	lastName := ctx.QueryParam(lastNameParam)
	userId := e.s.GetUserId(ctx)

	e.db.UpdateUserLastName(userId, lastName)
	ctx.JSON(http.StatusOK, true)
	return nil
}
