package mini_s3

import "github.com/labstack/echo"

type MiniS3Server struct {

}

func (this *MiniS3Server) ConnectMiniS3ToEcho(echo *echo.Echo) error {
	echo.GET("/api/arachne", this.Get)

	for serviceName, service := range this.services {
		echo.GET("/api/arachne/"+serviceName, service.listObjects)
		echo.GET("/api/arachne/"+serviceName+"/:id", service.handleObjectRequest)
		echo.GET("/api/arachne/"+serviceName+"/:id/snapshots", service.handleSnapshotListRequest)

	}
	return nil
}