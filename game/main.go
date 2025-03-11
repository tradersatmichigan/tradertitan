package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var gameId string

// add authentication to a handler
func LoginRequired(fn http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("gameId")

    if err == nil && cookie.Value == gameId {
      fn(w, r)
    } else {
      http.Error(w, "Please sign in!", http.StatusBadRequest)
    }
  }
}

func serveGame(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
    fmt.Fprintf(w, "Hello!") // serve file here
  } else {
    http.Error(w, "Bad method", http.StatusBadRequest)
  }
}

func handleJoin(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
    cookie, err := r.Cookie("gameId")

    if err == nil && cookie.Value == gameId {
      fmt.Fprintf(w, "Already signed in, please continue to the game")
    } else {
      http.ServeFile(w, r, "html/join.html")
    }

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

    mtx.Lock()
    if _, ok := users[username]; ok {
      http.Error(w, "Name already in use", http.StatusBadRequest)
    } else {
      users[username] = User{}

      http.SetCookie(w, &http.Cookie{
        Name: "gameId",
        Value: gameId,
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

      http.Redirect(w, r, "/", http.StatusOK)
    }
    mtx.Unlock()

  } else {
    http.Error(w, "Bad request method", http.StatusBadRequest)
  }
}

func main() {
  port := flag.String("port", "8080", "port to listen on")
  input := flag.String("file", "test.txt", "file path to load rounds from")

  flag.Parse()
  
  rounds := getRounds(*input)

  server := &http.Server{Addr: ":" + *port}

  // add http handlers
  http.HandleFunc("/join", handleJoin)
  http.HandleFunc("/", LoginRequired(serveGame))

  go func() {
    fmt.Println("server started on ", *port)
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
  }()

  RunGame(rounds)

  fmt.Println("Kill server?")
  waitForEnter() 

  ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
  defer cancel()

  if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	} else {
		fmt.Println("Server gracefully stopped.")
	}
}

func init() {
  // create a random number to verify login
  random, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
  if err != nil {
    panic(err)
  }
  gameId = random.String()
}

func getRounds(filepath string) []Round {
  file, err := os.Open(filepath)
  if err != nil {
    panic(err)
  }
  defer file.Close()

  var Rounds []Round

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    parts := strings.SplitN(line, " ", 2)

    number, err := strconv.Atoi(parts[0])
    if err != nil {
      panic(err)
    }

    Rounds = append(Rounds, Round{parts[1], number})
  }

  return Rounds
}

func waitForEnter() {
  bufio.NewReader(os.Stdin).ReadBytes('\n')
}
