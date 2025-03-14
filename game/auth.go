package main

import (
	"crypto/rand"
	"math"
	"math/big"
	"net/http"
	"time"
)

type ProtectedHandler = func(w http.ResponseWriter, r *http.Request, username string)

var IsLoggedIn func(*http.Request) bool

// require sign in for handler
var LoginRequired func(ProtectedHandler) http.HandlerFunc

// set cookies to login user
var Login ProtectedHandler


func init() {
  gameId, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))

  if err != nil {
    panic(err)
  }

  IsLoggedIn = func(r *http.Request) bool {
    if id, err := r.Cookie("gameId"); err != nil || 
       id.Value != gameId.String() {
      return false
    } else if _, err := r.Cookie("username"); err != nil {
      return false
    }
    return true
  }

  LoginRequired = func(hf ProtectedHandler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

      if IsLoggedIn(r) {
        username, _ := r.Cookie("username")
        hf(w, r, username.Value)
      } else {
        // http.Error(w, "Login Required", http.StatusBadRequest)
        http.Redirect(w, r, "/join", http.StatusSeeOther)
      }

    }
  }

  Login = func(w http.ResponseWriter, r *http.Request, username string) {
      http.SetCookie(w, &http.Cookie{
        Name: "gameId",
        Value: gameId.String(),
        Expires: time.Now().Add(3 * time.Hour),
        HttpOnly: true,
        Path: "/",
      })

      http.SetCookie(w, &http.Cookie{
        Name: "username",
        Value: username,
        Expires: time.Now().Add(3 * time.Hour),
        HttpOnly: true,
        Path: "/",
      })
  }
}
