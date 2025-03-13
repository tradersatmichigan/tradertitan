package main

import "sync"

var (
  users = make(map[Username] User)
  view = RegisterView
  mtx = sync.Mutex{}
  rooms []Room
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

// add catchup function
func GetUserChan(username Username) chan GameState {
  mtx.Lock()
  defer mtx.Unlock()

  if user, ok := users[username]; ok {
    return user.datachan
  } else {
    return nil
  }
}
