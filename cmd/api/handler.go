package main

import (
	"log"
	"net/http"
)

func (app *Config) TestSetup(w http.ResponseWriter, r *http.Request) {
	log.Println("Welcome")
}
