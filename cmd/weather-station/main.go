package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"weather-station/internal/cache"
	"weather-station/internal/fmi"
	"weather-station/internal/httpapi"
	"weather-station/internal/waterlevel"
)

func main() {
	c, err := loadConfig(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}
	service := waterlevel.NewService(fmi.NewClient(c.sourceURL), cache.Store{Dir: c.cacheDir}, c.now)
	service.Run(context.Background())
	mux := http.NewServeMux()
	mux.Handle("/api/water-level/hamina", httpapi.API{Service: service, Now: c.now})
	mux.Handle("/", http.FileServer(http.Dir(c.staticDir)))
	log.Printf("Keskeneräinen kehitysrunko: http://%s", c.listen)
	log.Fatal(http.ListenAndServe(c.listen, mux))
}
