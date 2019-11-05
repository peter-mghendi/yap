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
	try(godotenv.Load())
	try(model.InitDB(os.Getenv("DATABASE_URL")))
	try(handler.InitRBAC())
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	port = os.Getenv("PORT")
}

/* TODO
1. Forum
2. Relational data
3. User actions
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

	u := e.Group("/users")
	u.GET("", handler.GetUsers)
	u.GET("/:id", handler.GetUserByID)
	u.POST("/join", handler.JoinUser)
	u.POST("/auth", handler.AuthUser)
	// u.GET("/:id/blog/posts", handler.GetUserBlogPosts) // TODO
	u.GET("/:id/blog/reactions", handler.GetUserBlogReactions)
	// u.GET("/:id/forum/reactions", handler.GetUserForumReactions) // TODO
	// u.GET("/:id/forum/questions", handler.GetUserForumQuestions) // TODO
	// u.GET("/:id/forum/responses", handler.GetUserForumResponses) // TODO

	uAuth := u.Group("/restricted")
	uAuth.Use(middleware.JWTWithConfig(jwtConfig))
	uAuth.PUT("/me/update", handler.UpdateUser)
	uAuth.PUT("/me/change", handler.UpdatePass)
	uAuth.PUT("/:id/assign", handler.AssignUser)
	uAuth.DELETE("/:id/delete", handler.DeleteUser)

	/* TODO
	1. RSS Feeds
	2. All posts
	3. Search
	*/
	blog := e.Group("/blog")
	// blog.GET("/feed", handler.GetFeed) // TODO

	p := blog.Group("/posts")
	pAuth := p.Group("/:id")
	pAuth.Use(middleware.JWTWithConfig(jwtConfig))
	pAuth.DELETE("/delete", handler.DeletePost)
	pAuth.PUT("/publish", handler.PublishPost)
	pAuth.PUT("/retract", handler.RetractPost)

	pr := p.Group("/:id/reactions")

	pr.GET("", handler.GetBlogPostReactions)

	prAuth := pr.Group("/create")
	prAuth.Use(middleware.JWTWithConfig(jwtConfig))
	prAuth.POST("", handler.CreateBlogReaction)

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

	br := blog.Group("/reactions")
	br.GET("/:id", handler.GetBlogReactionByID)

	brAuth := br.Group("")
	brAuth.Use(middleware.JWTWithConfig(jwtConfig))
	brAuth.PUT("/:id/update", handler.UpdateBlogReaction)
	brAuth.DELETE("/:id/delete", handler.DeleteBlogReaction)

	/* TODO
	1. Password Reset
	*/

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

// try handles top-level errors
func try(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
