package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"weather-station/internal/cache"
	"weather-station/internal/fmi"
	"weather-station/internal/httpapi"
	"weather-station/internal/waterlevel"
)

func handler(staticDir string, service *waterlevel.Service, now func() time.Time) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/api/water-level/hamina", httpapi.API{Service: service, Now: now})
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))
	return mux
}

func main() {
	c, err := loadConfig(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}
	service := waterlevel.NewService(fmi.NewClient(c.sourceURL), cache.Store{Dir: c.cacheDir}, c.now)
	service.Run(context.Background())
	log.Printf("Keskeneräinen kehitysrunko: http://%s", c.listen)
	log.Fatal(http.ListenAndServe(c.listen, handler(c.staticDir, service, c.now)))
}
