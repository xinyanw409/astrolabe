package main 

import (
	//"net/http"
	"github.com/labstack/echo"
	"github.com/vmware/arachne/rest_api"
)

func main() {
	arachneRestAPI:= rest_api.NewArachne()
	e := echo.New()
	err := arachneRestAPI.ConnectArachneAPIToEcho(e)
	if (err != nil) {
		e.Logger.Fatal(err)
	}
	e.Logger.Fatal(e.Start(":1323"))
}

