package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	m "github.com/l3njo/yap/models"
	u "github.com/l3njo/yap/utils"

	"github.com/joho/godotenv"
)

var (
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

	u.Try(godotenv.Load())
	u.Try(m.InitDB(os.Getenv("DATABASE_URL")))
	port = os.Getenv("PORT")
}

func main() {
}
