package main

import (
	"MIS/pkg/settings"
	"MIS/routers"
	"log"
)

func main() {
	engine := routers.RoutesController()
	err := engine.Run(":" + settings.ServerSettings.HttpPort)
	if err != nil {
		log.Fatalf("Fail to Run engine:%v", err)
	}
}
