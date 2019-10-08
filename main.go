package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/l3njo/yap/db"
	"github.com/l3njo/yap/handler"
	"github.com/l3njo/yap/model"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	e         *echo.Echo
	port      string
	jwtSecret []byte
	signals   chan os.Signal
)

func cleanup() {
	log.Println("Shutting down server.")
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
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	port = os.Getenv("PORT")
}

func main() {
	jwtConfig := middleware.JWTConfig{
		Claims:     &handler.JwtCustomClaims{},
		SigningKey: jwtSecret,
	}

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	u := e.Group("/users")
	u.GET("", handler.GetUsers)
	u.GET("/:id", handler.GetUserByID)
	u.POST("/join", handler.JoinUser)
	u.POST("/auth", handler.AuthUser)

	uAuth := u.Group("/:id")
	uAuth.Use(middleware.JWTWithConfig(jwtConfig))
	uAuth.PUT("/update", handler.UpdateUser)
	uAuth.PUT("/assign", handler.AssignUser)
	uAuth.DELETE("/delete", handler.DeleteUser)

	p := e.Group("/posts")
	pAuth := p.Group("/:id")
	pAuth.Use(middleware.JWTWithConfig(jwtConfig))
	pAuth.DELETE("/delete", handler.DeletePost)
	pAuth.PUT("/publish", handler.PublishPost)
	pAuth.PUT("/retract", handler.RetractPost)

	a := p.Group("/articles")
	a.GET("/public", handler.GetPublicArticles)
	a.GET("/public/:id", handler.GetPublicArticleByID)

	aAuth := a.Group("")
	aAuth.Use(middleware.JWTWithConfig(jwtConfig))
	aAuth.GET("", handler.GetArticles)
	aAuth.GET("/:id", handler.GetArticleByID)
	aAuth.POST("/create", handler.CreateArticle)
	aAuth.PUT("/:id/update", handler.UpdateArticle)

	g := p.Group("/galleries")
	g.GET("/public", handler.GetPublicGalleries)
	g.GET("/public/:id", handler.GetPublicGalleryByID)

	gAuth := g.Group("")
	gAuth.Use(middleware.JWTWithConfig(jwtConfig))
	gAuth.GET("", handler.GetGalleries)
	gAuth.GET("/:id", handler.GetGalleryByID)
	gAuth.POST("/create", handler.CreateGallery)
	gAuth.PUT("/:id/update", handler.UpdateGallery)

	f := p.Group("/flickers")
	f.GET("/public", handler.GetPublicFlickers)
	f.GET("/public/:id", handler.GetPublicFlickerByID)

	fAuth := f.Group("")
	fAuth.Use(middleware.JWTWithConfig(jwtConfig))
	fAuth.GET("", handler.GetFlickers)
	fAuth.GET("/:id", handler.GetFlickerByID)
	fAuth.POST("/create", handler.CreateFlicker)
	fAuth.PUT("/:id/update", handler.UpdateFlicker)

	r := e.Group("/reactions")
	r.GET("", handler.GetReactions)
	r.GET("/:id", handler.GetReactionByID)

	rAuth := r.Group("")
	rAuth.Use(middleware.JWTWithConfig(jwtConfig))
	rAuth.POST("/create", handler.CreateReaction)
	rAuth.PUT("/:id/update", handler.UpdateReaction)
	rAuth.DELETE("/:id/delete", handler.DeleteReaction)

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
		log.Fatalln(err)
	}
}
