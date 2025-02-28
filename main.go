package main

import (
	"fmt"
	"log"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
    http.ServeFile(w, r, "join.html")
  } else if r.Method == http.MethodPost {
    if r.ParseForm() != nil {
      log.Fatal("Form parse fail")
    }
    fmt.Fprintf(w, "%s", r.Form.Get("username"))
  } else {
    http.Error(w, "Bad Request", http.StatusBadRequest)
  }
}

func main() {
  http.HandleFunc("/join/", loginHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}
