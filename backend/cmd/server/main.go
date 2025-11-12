package server

import (
	"log"
)

func main() {
	app := SetupApp()

	log.Println("listening on port", app.Config.Port)
	if err := app.Router.Run(":" + app.Config.Port); err != nil {
		log.Fatal(err)
	}
}
