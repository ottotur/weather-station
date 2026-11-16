package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Keskeneräinen kehitysrunko: http://127.0.0.1:8080")
	if err := http.ListenAndServe("127.0.0.1:8080", http.FileServer(http.Dir("web/dist"))); err != nil {
		panic(err)
	}
}
