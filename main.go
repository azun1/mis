package main

import (
	"MIS/routers"
)

func main() {
	router := routers.InitRouter()

	router.Run()
}
