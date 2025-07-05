package main

/*
App entry point
*/
import (
	"web-page-analyzer/router"
)

func main() {

	r := router.SetupRouter()
	r.Run(":8080") // start server woith http://localhost:8080

}
