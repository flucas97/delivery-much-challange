package main

import (
	"github.com/flucas97/delivery-much-challange/internal/config/router"
	"github.com/flucas97/delivery-much-challange/tools/loggertools"
)

func main() {
	loggertools.Info("Hello Delivery Much Team, welcome :)")
	router.StartRouter()
}
