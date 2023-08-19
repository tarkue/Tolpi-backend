package app

import (
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
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

	a.echo.Use(middleware.Logger())
	a.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowCredentials: true,
	}))
	a.echo.Use(a.middleware.Authorization)

	a.echo.POST("/subscribe", a.e.Subscribe)
	a.echo.POST("/unsubscribe", a.e.Unsubscribe)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
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
