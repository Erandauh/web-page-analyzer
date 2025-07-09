package main

/*
App entry point
*/
import (
	"web-page-analyzer/logger"
	"web-page-analyzer/router"

	"github.com/sirupsen/logrus"
)

func main() {

	logger.InitLogger()

	logrus.Info("Starting server...")
	r := router.SetupRouter()
	r.Run(":8080") // start server woith http://localhost:8080

}
