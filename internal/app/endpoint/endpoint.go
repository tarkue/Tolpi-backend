package endpoint

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tarkue/tolpi-backend/internal/app/graph/model"
)

type DataBase interface {
	FindUserById(userID string) *model.User
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

func (e *Endpoint) Subscribe(ctx echo.Context) error {
	if !ctx.QueryParams().Has(subIdParam) {
		ctx.String(http.StatusOK, notFoundMessage+subIdParam)
		return nil
	}

	userId := e.s.GetUserId(ctx)
	subId := ctx.QueryParam(subIdParam)
	user := e.db.FindUserById(subId)

	if e.s.Contains(user.TrackerList, userId) {
		ctx.JSON(http.StatusOK, true)
		return nil
	}

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

	if !e.s.Contains(user.TrackerList, userId) {
		ctx.JSON(http.StatusOK, true)
		return nil
	}

	NewSlice := e.s.RemoveIndex(user.TrackerList, e.s.IndexOf(userId, user.TrackerList))
	e.db.UpdateUserTrackers(unsubId, NewSlice)

	ctx.JSON(http.StatusOK, true)
	return nil
}
