package main

import "sync"

var (
  mtx = sync.Mutex{}
  users = make(map[Username] User)
)

