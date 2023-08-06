package app

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/tarkue/tolpi-backend/config"
	"github.com/tarkue/tolpi-backend/internal/app/database"
	"github.com/tarkue/tolpi-backend/internal/app/endpoint"
	"github.com/tarkue/tolpi-backend/internal/app/graph"
)

type App struct {
	e    *endpoint.Endpoint
	db   *database.DB
	echo *echo.Echo
}

func New() (*App, error) {
	a := &App{}

	DataBase := database.New()

	a.db = DataBase

	a.e = endpoint.New(a.db)

	a.echo = echo.New()

	a.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	playground := playground.Handler("GraphQL playground", "/query")

	a.echo.GET("/", func(c echo.Context) error {
		playground.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	a.echo.GET("/getStatus", a.e.GetStatus)
	a.echo.GET("/getCountry", a.e.GetCountry)

	a.echo.POST("/query", func(c echo.Context) error {
		srv.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	return a, nil

}

func (a *App) Run() error {

	err := a.echo.Start(config.ServerPort)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("server running...")
	return nil
}
