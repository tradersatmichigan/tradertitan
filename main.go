package main

import (
	"net/http"
)

func handleJoin(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {

  } else if r.Method == http.MethodPut {

  } else {

  }
}

func main() {
  http.HandleFunc("/join", handleJoin)
}
