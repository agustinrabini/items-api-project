package main

import (
	"github.com/agustinrabini/items-api-project/src/main/api/app"
	"github.com/agustinrabini/items-api-project/src/main/api/config"
)

// @title API Items
// @version 1.0
// @description This is a api items.
// @termsOfService http://swagger.io/terms/

// @contact.name Agustin Rabini
// @contact.url http://www.swagger.io/support
// @contact.email agustinrabini99@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	config.Load()
	app.Start()
}
