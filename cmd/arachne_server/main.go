package main

import (
	"flag"
	//"github.com/labstack/gommon/log"
	"log"
	//"net/http"
	"github.com/labstack/echo"
	"github.com/vmware/arachne/pkg/rest_api"
)

func main() {
	confDirStr := flag.String("confDir", "", "Configuration directory")
	apiPortStr := flag.String("apiPort", "1323", "REST API port")
	flag.Parse()
	if *confDirStr == "" {
		log.Println("confDir is not defined")
		flag.Usage()
		return
	}
	arachneRestAPI := rest_api.NewArachne(*confDirStr)
	e := echo.New()
	err := arachneRestAPI.ConnectArachneAPIToEcho(e)
	if (err != nil) {
		e.Logger.Fatal(err)
	}
	e.Logger.Fatal(e.Start(":" + *apiPortStr))
}

