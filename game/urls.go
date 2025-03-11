package main

import "net/http"

func init() {
  http.HandleFunc("/join", HandleLogin)
  http.HandleFunc("/", LoginRequired(ServeGame))
}
