package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

// buffer ts
func GetStream(w http.ResponseWriter, r *http.Request, username string) {
  if r.Method == http.MethodGet {
    datachan := GetUserChan(username)

    if datachan != nil {
      w.Header().Set("Content-Type", "text/event-stream")
      w.Header().Set("Cache-Control", "no-cache")
      w.Header().Set("Connection", "keep-alive")

      flusher, ok := w.(http.Flusher)
      if !ok {
        http.Error(w, "Streaming not supported", http.StatusInternalServerError)
        return
      }

      for {
        select {
        case <- r.Context().Done():
          return
        case state, ok := <- datachan:
          if !ok {
            return
          }

          data, err := json.Marshal(state)
          if err != nil {
            http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
            return
          }

          fmt.Fprintf(w, "data: %s\n\n", data)
          flusher.Flush()
        }
      }
    } else {
      http.Error(w, "Bad username", http.StatusBadRequest)
    }
  } else {
    http.Error(w, "Bad request method", http.StatusBadRequest)
  }
}

func PostMake(w http.ResponseWriter, r *http.Request, username string) {
  if r.Method == http.MethodPost {

    if r.ParseForm() != nil {
      http.Error(w, "Form error", http.StatusBadRequest)
      return
    }

    formVal := r.Form.Get("value")

    width, err := strconv.ParseUint(formVal, 10, 64)

    if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    }

    mtx.Lock()
    if user, ok := users[username]; ok && view == MakeView {
      room := &rooms[user.room]
      
      if room.width > width {
        room.username = username
        room.width = width
      }
    }
    mtx.Unlock()

    w.WriteHeader(http.StatusOK) 
    w.Write([]byte("Success"))

  } else {
    http.Error(w, "Bad request method", http.StatusBadRequest)
  }
}

// check negative?
func PostCenter(w http.ResponseWriter, r *http.Request, username string) {
  if r.Method == http.MethodPost && r.ParseForm() == nil {

    if r.ParseForm() != nil {
      http.Error(w, "Form error", http.StatusBadRequest)
      return
    }

    formVal := r.Form.Get("value")

    center, err := strconv.ParseUint(formVal, 10, 64)

    if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    }

    mtx.Lock()
    if user, ok := users[username]; ok && view == CenterView {
      room := &rooms[user.room]
      
      if room.username == username {
        room.center = center
      }
    }
    mtx.Unlock()

    w.WriteHeader(http.StatusOK) 
    w.Write([]byte("Success"))

  } else {
    http.Error(w, "Bad request method", http.StatusBadRequest)
  }
}

func PostTrade(w http.ResponseWriter, r *http.Request, username string) {
  if r.Method == http.MethodPost {

    if r.ParseForm() != nil {
      http.Error(w, "Form error", http.StatusBadRequest)
      return
    }

    formVal := r.Form.Get("value")

    mtx.Lock()
    if user, ok := users[username]; ok && view == TradeView && rooms[user.room].username != "" {

      if formVal == "buy" {
        user.side = Long
      } else if formVal == "sell" {
        user.side = Short
      }

    }
    mtx.Unlock()

    w.WriteHeader(http.StatusOK) 
    w.Write([]byte("Success"))

  } else {
    http.Error(w, "Bad request method", http.StatusBadRequest)
  }
}
