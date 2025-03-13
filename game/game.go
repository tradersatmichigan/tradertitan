package main

import "sync"

var (
  users = make(map[Username] User)
  view = RegisterView
  mtx = sync.Mutex{}
)
func Register(username Username) bool {
  mtx.Lock()
  defer mtx.Unlock()

  if view != RegisterView {
    return false
  } else if _, ok := users[username]; ok {
    return false
  } 

  users[username] = User{}

  return true
}

func Make(username Username, price uint) {
  mtx.Lock()
  defer mtx.Unlock()
}
