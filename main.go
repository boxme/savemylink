package main

import (
	"flag"
	"net/http"
	"savemylink/auth"
	"savemylink/database"
	"strconv"
)

type App struct {
	AuthHandler *auth.AuthHandler
}

func main() {
	boolPtr := flag.Bool("prod", false, "Provide true for this flag in production")
	flag.Parse()

	cfg := LoadConfig(*boolPtr)
	dbService := database.NewDbService(cfg.Database)
	a := &App{
		AuthHandler: auth.NewAuthHandler(dbService),
	}

	defer dbService.Close()

	port := ":" + strconv.Itoa(cfg.Port)
	http.ListenAndServe(port, a)
}

func (a *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	a.AuthHandler.ServeHTTP(res, req)
}
