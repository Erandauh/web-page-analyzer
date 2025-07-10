package main

/*
App entry point
*/
import (
	"web-page-analyzer/logger"

	"github.com/sirupsen/logrus"
)

func main() {

	logger.InitLogger()

	logrus.Info("Starting server...")
	r := InitializeRouter()
	r.Run(":8080") // start server woith http://localhost:8080

}
