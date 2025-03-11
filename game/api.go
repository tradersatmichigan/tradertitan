package main

import (
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
  if IsLoggedIn(r) {
    http.Error(w, "Already signed in, please continue to the game", http.StatusBadRequest)
  } else if r.Method == http.MethodGet {
    http.ServeFile(w, r, "html/join.html")
  } else if r.Method == http.MethodPost {

    if err := r.ParseForm(); err != nil {
      http.Error(w, "Error parsing form", http.StatusBadRequest)
      return
    }

		username := r.FormValue("username")
		if username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}

    if Register(username) {
      Login(w, r, username)
      http.Redirect(w, r, "/", http.StatusSeeOther)
    } else {
      http.Error(w, "Name already in use", http.StatusBadRequest)
    }
    
  } else {
    http.Error(w, "Bad request method", http.StatusBadRequest)
  }
}

func ServeGame(w http.ResponseWriter, r *http.Request, _ string) { // also return current state
  if r.Method == http.MethodGet {
    http.ServeFile(w, r, "html/game.html")
  } else {
    http.Error(w, "Bad request method", http.StatusBadRequest)
  }
}

func GetStream(w http.ResponseWriter, r *http.Request, username string) {
  if r.Method == http.MethodGet {

  } else {
    http.Error(w, "Bad request method", http.StatusBadRequest)
  }
}

func PostWidth(w http.ResponseWriter, r *http.Request, username string) {
  if r.Method == http.MethodPost {

  } else {
    http.Error(w, "Bad request method", http.StatusBadRequest)
  }
}

func PostCenter(w http.ResponseWriter, r *http.Request, username string) {
  if r.Method == http.MethodPost {
    if r.ParseForm() == nil {

    }
  } else {
    http.Error(w, "Bad request method", http.StatusBadRequest)
  }
}

func PostTrade(w http.ResponseWriter, r *http.Request, username string) {
  if r.Method == http.MethodPost {

  } else {
    http.Error(w, "Bad request method", http.StatusBadRequest)
  }
}
