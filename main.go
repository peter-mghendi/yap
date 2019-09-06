package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/l3njo/yap-api/db"
	"github.com/l3njo/yap-api/handler"
	"github.com/l3njo/yap-api/model"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	e       *echo.Echo
	port    string
	signals chan os.Signal
)

func cleanup() {
	e.Logger.Info("Shutting down server.")
	db.DB.Close()
}

func init() {
	signals = make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		cleanup()
		os.Exit(1)
	}()

	e = echo.New()
	Try(godotenv.Load())
	Try(model.InitDB(os.Getenv("DATABASE_URL")))
	port = os.Getenv("PORT")
}

func main() {
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	u := e.Group("/users") // TODO
	u.GET("", handler.AppController)

	p := e.Group("/posts")
	p.DELETE("/:id/delete", handler.DeletePost)
	p.PUT("/:id/publish", handler.PublishPost)
	p.PUT("/:id/retract", handler.RetractPost)

	a := p.Group("/articles")
	a.GET("", handler.GetArticles)
	a.GET("/:id", handler.GetArticleByID)
	a.POST("/create", handler.CreateArticle)
	a.PUT("/:id/update", handler.UpdateArticle)

	g := p.Group("/galleries")
	g.GET("", handler.GetGalleries)
	g.GET("/:id", handler.GetGalleryByID)
	g.POST("/create", handler.CreateGallery)
	g.PUT("/:id/update", handler.UpdateGallery)

	f := p.Group("/flickers")
	f.GET("", handler.GetFlickers)
	f.GET("/:id", handler.GetFlickerByID)
	f.POST("/create", handler.CreateFlicker)
	f.PUT("/:id/update", handler.UpdateFlicker)

	r := e.Group("/reactions")
	r.GET("", handler.GetReactions)
	r.GET("/:id", handler.GetReactionByID)
	r.POST("/create", handler.CreateReaction)
	r.PUT("/:id/update", handler.UpdateReaction)
	r.DELETE("/:id/delete", handler.DeleteReaction)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		resp := handler.Response{
			Message: http.StatusText(code),
		}

		c.JSON(code, resp)
		c.Logger().Error(err)
	}

	e.GET("/", handler.AppController)
	e.Logger.Fatal(e.Start(":" + port))
}

// Try handles top-level errors
func Try(err error) {
	if err != nil {
		e.Logger.Fatal(err)
	}
}
