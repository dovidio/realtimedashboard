package main

import (
	"realtimedashboard/appdownload"
	"realtimedashboard/database"
)

func main() {
	database.SetupDatabase()
	appdownload.SetupRoutes()
}
