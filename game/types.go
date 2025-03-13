package main

type Username = string

type User struct {
  room uint
  totalPnl int
  totalPlace uint
}

type View = int

const (
  RegisterView View = iota
  MakeView
  CenterView
  TradeView
) // add leader view?

type GameState struct {
  view View
}

type Round struct {
  Market string
  TrueValue int
}
