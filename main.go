package main

import (
	"fmt"
	"log"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
    fmt.Fprintf(w, "Got request")
  } else if r.Method == http.MethodPost {

  } else {

  }
}

func main() {
  http.HandleFunc("/join/", loginHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}
