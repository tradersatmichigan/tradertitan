package main

type Username = string

type User struct {
  room uint
  totalPnl int
  totalPlace uint
}

type Round struct {
  Market string
  TrueValue int
}

type View = int
