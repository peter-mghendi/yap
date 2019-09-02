package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	m "github.com/l3njo/yap/models"
	c "github.com/l3njo/yap/controllers"

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
	log.Println("Shutting down server.")
	m.DB.Close()
}

func init() {
	signals = make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		cleanup()
		os.Exit(1)
	}()

	Try(godotenv.Load())
	Try(m.InitDB(os.Getenv("DATABASE_URL")))
	port = os.Getenv("PORT")
}

func main() {
	e = echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.GET("/", c.AppController)

	e.Logger.Fatal(e.Start(":" + port))
}

// Try handles top-level errors
func Try(err error) {
	if err != nil {
		e.Logger.Fatal(err)
	}
}
