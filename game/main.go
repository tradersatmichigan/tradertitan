package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)


func main() {
  port := flag.String("port", "8080", "port to listen on")
  input := flag.String("file", "test.txt", "file path to load rounds from")

  flag.Parse()
  
  rounds := getRounds(*input)
  server := &http.Server{Addr: ":" + *port}

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
