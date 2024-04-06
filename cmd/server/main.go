package main

import (
	"fmt"
	"github.com/bigusef/texorbit/internal/router"
	"github.com/bigusef/texorbit/pkg/config"
	"log"
	"net/http"
)

func main() {
	setting := config.NewSetting()
	handler := router.New(setting)

	// start application server
	log.Println("[Restfull server] Start Store application")
	log.Println(fmt.Sprintf("Serving starting on port :%s", setting.Port))

	serverPort := fmt.Sprintf(":%s", setting.Port)
	if err := http.ListenAndServe(serverPort, handler); err != nil {
		log.Fatal(err)
	}
}
