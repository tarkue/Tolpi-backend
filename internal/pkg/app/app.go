package app

import (
	"log"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/tarkue/tolpi-backend/config"
	"github.com/tarkue/tolpi-backend/internal/app/database"
	"github.com/tarkue/tolpi-backend/internal/app/endpoint"
	"github.com/tarkue/tolpi-backend/internal/app/graph"
	"github.com/tarkue/tolpi-backend/internal/app/middlewares"
	"github.com/tarkue/tolpi-backend/internal/app/service"
)

type App struct {
	s          *service.Service
	e          *endpoint.Endpoint
	db         *database.DB
	middleware *middlewares.Middlewares
	echo       *echo.Echo
}

func New() (*App, error) {
	a := &App{}

	a.s = service.New()
	DataBase := database.New()

	a.db = DataBase

	a.middleware = middlewares.New(a.s)

	a.e = endpoint.New(a.db, a.s)

	a.echo = echo.New()

	a.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}), a.middleware.Authorization)

	a.echo.GET("/getCountry", a.e.GetCountry)

	a.echo.POST("/subscribe", a.e.Subscribe)
	a.echo.POST("/unsubscribe", a.e.Unsubscribe)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})

	a.echo.GET("/", echo.WrapHandler(playground.Handler("GraphQL playground", "/query")))
	a.echo.Any("/query", echo.WrapHandler(srv))

	return a, nil

}

func (a *App) Run() error {

	err := a.echo.Start(":" + config.ServerPort)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("server running...")
	return nil
}
