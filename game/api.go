package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if IsLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
			http.Error(w, "Name already in use or game has already started",
				http.StatusBadRequest)
		}

	} else {
		http.Error(w, "Bad request method", http.StatusBadRequest)
	}
}

func ServeGame(w http.ResponseWriter, r *http.Request, _ string) { // also return current state
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "static/index.html")
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
				case <-r.Context().Done():
					return
				case state, ok := <-datachan:
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

		formVal := r.FormValue("value")

		width, err := strconv.ParseFloat(formVal, 64)

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mtx.Lock()
		if user, ok := users[username]; ok && view == MakeView {
			room := &rooms[user.room]

			if room.Width > width || room.Username == "" {
				room.Username = username
				room.Width = width

				for otheruser := range users {
					if user.room == users[otheruser].room {
						PushUserState(otheruser)
					}
				}
			}
		}
		mtx.Unlock()

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Bad request method", http.StatusBadRequest)
	}
}

// check negative?
func PostCenter(w http.ResponseWriter, r *http.Request, username string) {
	if r.Method == http.MethodPost {

		formVal := r.FormValue("value")

		center, err := strconv.ParseFloat(formVal, 64)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mtx.Lock()
		if user, ok := users[username]; ok && view == CenterView {
			room := &rooms[user.room]

			if room.Username == username {
				room.Center = center
			}

			PushUserState(username)
		}
		mtx.Unlock()

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Bad request method", http.StatusBadRequest)
	}
}

func PostTrade(w http.ResponseWriter, r *http.Request, username string) {
	if r.Method == http.MethodPost {

		formVal := r.FormValue("value")

		mtx.Lock()
		if user, ok := users[username]; ok && view == TradeView && rooms[user.room].Username != "" {

			if formVal == "buy" {
				user.side = Long
			} else if formVal == "sell" {
				user.side = Short
			}

		}
		mtx.Unlock()

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Bad request method", http.StatusBadRequest)
	}
}
