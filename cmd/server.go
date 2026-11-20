package cmd

import (
	"log"
	"net/http"
	"gin-spa-app-framework-back/initialize"
	"gin-spa-app-framework-back/internal/routes"
)

func startWebServer() {
	r := routes.Router
	webServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("webserver listen error: %s\n", err)
	}

}

func RunServer() {
	initialize.InitSystem()
	startWebServer()
}
