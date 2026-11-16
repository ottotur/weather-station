package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	c, err := loadConfig(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Keskeneräinen kehitysrunko: http://%s", c.listen)
	log.Fatal(http.ListenAndServe(c.listen, http.FileServer(http.Dir(c.staticDir))))
}
