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
	_ = db.DB.Close()
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
	try(godotenv.Load())
	try(model.InitDB(os.Getenv("DATABASE_URL")))
	try(handler.InitRBAC())
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	port = os.Getenv("PORT")
}

/* TODO
1. Podcasts.
2. Relational data (posts, reactions).
3. User data (separately per post type).
4. Search (separately per post type).
5. Password Reset.
*/
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

	// PATH /users
	u := e.Group("/users")
	u.GET("", handler.GetUsers)
	u.GET("/:id", handler.GetUserByID)
	u.POST("/join", handler.JoinUser)
	u.POST("/auth", handler.AuthUser)
	u.GET("/:id/posts/articles", handler.GetUserPublicArticles) // TODO Relations
	u.GET("/:id/posts/galleries", handler.GetUserPublicGalleries) // TODO Relations
	u.GET("/:id/posts/flickers", handler.GetUserPublicFlickers) // TODO Relations
	u.GET("/:id/reactions", handler.GetUserReactions) // TODO Relations

	// PATH /users/restricted
	uAuth := u.Group("/restricted")
	uAuth.Use(middleware.JWTWithConfig(jwtConfig))
	// uAuth.GET("/:id/posts/articles", handler.GetUserArticles) // TODO Relations
	// uAuth.GET("/:id/posts/galleries", handler.GetUserGalleries) // TODO Relations
	// uAuth.GET("/:id/posts/flickers", handler.GetUserFlickers) // TODO Relations
	uAuth.PUT("/me/update", handler.UpdateUser)
	uAuth.PUT("/me/change", handler.UpdatePass)
	uAuth.PUT("/:id/assign", handler.AssignUser)
	uAuth.DELETE("/:id/delete", handler.DeleteUser)

	// PATH /posts
	p := e.Group("/posts")
	pAuth := p.Group("/:id")
	pAuth.Use(middleware.JWTWithConfig(jwtConfig))
	pAuth.DELETE("/delete", handler.DeletePost)
	pAuth.PUT("/publish", handler.PublishPost)
	pAuth.PUT("/retract", handler.RetractPost)

	// PATH /posts/:id/reactions
	pr := p.Group("/:id/reactions") // TODO Relations
	pr.GET("", handler.GetPostReactions) // TODO Relations
	pr.GET("/:reaction", handler.GetPostReactionByID) // TODO Relations

	// PATH /posts/:id/reactions/restricted
	prAuth := pr.Group("/restricted")
	prAuth.Use(middleware.JWTWithConfig(jwtConfig))
	prAuth.POST("/create", handler.CreateReaction)
	prAuth.PUT("/:reaction/update", handler.UpdateReaction)
	prAuth.DELETE("/:reaction/delete", handler.DeleteReaction)

	// PATH /posts/articles
	a := p.Group("/articles")
	a.GET("/public", handler.GetPublicArticles)
	a.GET("/public/:id", handler.GetPublicArticleByID)

	aAuth := a.Group("")
	aAuth.Use(middleware.JWTWithConfig(jwtConfig))
	aAuth.GET("", handler.GetArticles)
	aAuth.GET("/:id", handler.GetArticleByID)
	aAuth.POST("/create", handler.CreateArticle)
	aAuth.PUT("/:id/update", handler.UpdateArticle)
	aAuth.PUT("/:id/transfer", handler.TransferArticle)

	// PATH /posts/galleries
	g := p.Group("/galleries")
	g.GET("/public", handler.GetPublicGalleries)
	g.GET("/public/:id", handler.GetPublicGalleryByID)

	gAuth := g.Group("")
	gAuth.Use(middleware.JWTWithConfig(jwtConfig))
	gAuth.GET("", handler.GetGalleries)
	gAuth.GET("/:id", handler.GetGalleryByID)
	gAuth.POST("/create", handler.CreateGallery)
	gAuth.PUT("/:id/update", handler.UpdateGallery)
	gAuth.PUT("/:id/transfer", handler.TransferGallery)

	// PATH /posts/flickers
	f := p.Group("/flickers")
	f.GET("/public", handler.GetPublicFlickers)
	f.GET("/public/:id", handler.GetPublicFlickerByID)

	fAuth := f.Group("")
	fAuth.Use(middleware.JWTWithConfig(jwtConfig))
	fAuth.GET("", handler.GetFlickers)
	fAuth.GET("/:id", handler.GetFlickerByID)
	fAuth.POST("/create", handler.CreateFlicker)
	fAuth.PUT("/:id/update", handler.UpdateFlicker)
	fAuth.PUT("/:id/transfer", handler.TransferFlicker)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		resp := handler.Response{
			Message: http.StatusText(code),
		}

		_ = c.JSON(code, resp)
		c.Logger().Error(err)
	}

	e.GET("/", handler.AppController)
	e.Logger.Fatal(e.Start(":" + port))
}

// try handles top-level errors
func try(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
