package main

import (
	"sync"
	"time"
)

var (
  mtx = sync.Mutex{}
  users = make(map[Username] User)
)

func SleepAndFree(time.Duration) {
  mtx.Unlock()
  time.Sleep()
}

func RunGame(Rounds []Round) {
  mtx.Lock()
  defer mtx.Unlock()

  for _, round := range Rounds {

  }

}
