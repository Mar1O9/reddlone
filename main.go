package main

import (
	"github.com/Mar1O9/reddlone/api"
	"github.com/Mar1O9/reddlone/app"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	a := &app.App{}

	serv := api.InitServer()

	a.RunApp(serv)
}
