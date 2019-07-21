package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	//"github.com/labstack/gommon/log"
	"log"
	//"net/http"
	"github.com/labstack/echo"
	"github.com/vmware/arachne/pkg/server"
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
	apiPort, err := strconv.Atoi(*apiPortStr)
	if err != nil {
		fmt.Errorf("apiPort %s is not an integer", *apiPortStr)
		os.Exit(1)
	}
	arachneRestAPI := server.NewArachne(*confDirStr, apiPort)
	e := echo.New()
	err = arachneRestAPI.ConnectArachneAPIToEcho(e)
	if err != nil {
		e.Logger.Fatal(err)
	}
	err = arachneRestAPI.ConnectMiniS3ToEcho(e)
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Fatal(e.Start(":" + *apiPortStr))
}

