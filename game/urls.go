package main

import "net/http"

func init() {
	http.HandleFunc("/join", HandleLogin)
	http.HandleFunc("/", LoginRequired(ServeGame))
	http.HandleFunc("/event", LoginRequired(GetStream))

	http.HandleFunc("/make", LoginRequired(PostMake))
	http.HandleFunc("/center", LoginRequired(PostCenter))
	http.HandleFunc("/trade", LoginRequired(PostTrade))
}
