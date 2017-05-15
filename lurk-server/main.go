package main

import (
	"log"
	"net/http"

	"../starstore"
)

func main() {
	dbfile := "./registration.db"

	err := starstore.Init(dbfile)
	if err != nil {
		log.Fatalf("%s: %s", dbfile, err)
	}

	log.Fatal(http.ListenAndServe(":8080", NewRouter()))
}
