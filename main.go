package main

import (
	"log"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
    http.ServeFile(w, r, "join.html")
  } else if r.Method == http.MethodPost {

  } else {
  }
}

func main() {
  http.HandleFunc("/join/", loginHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}
