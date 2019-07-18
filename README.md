Welcome to Arachne

To build

go build

To install

cd arachne_server
go install

To run the Arachne server
$GOPATH/bin/arachne_server -confDir=<your configuration dir> [-port=<desired port number>]

The default port for the server is 1323

Access via <ip>:/api/arachne

VMWare IVDs (Improved Virtual Disks aka First Class Disks/FCDs) are supported

Configure with a configuration directory that you specify with -confDir.  
To configure IVD a file named ivd.pe.json needs to be created

{
	"vcHost":"<host name/IP>",
	"insecureVC":"<Y unless your server has a proper certificate>",
	"vcUser":"<VC login>",
	"vcPassword":"<VC password>"
}