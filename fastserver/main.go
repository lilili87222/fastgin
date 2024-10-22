package main

import (
	"fastgin/boost"
)

// @title Go Web fastgin API
// @version 1.0
// @description This is a sample server for a Go web mini project.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 192.168.123.214:8088
// @BasePath /
func main() {
	boost.StartWebService()
}
