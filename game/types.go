package main

type Username = string

type User struct {
  room uint
  totalPnl int
  totalPlace uint

  side Side
  datachan chan GameState
}

type View = int

const (
  RegisterView View = iota
  MakeView
  CenterView
  TradeView
) // add leader view?

type Side = int

const (
  Long Side = iota
  Short
  None
)

type GameState struct {
  view View
  room Room
}

type Round struct {
  Market string
  TrueValue int
}

type Room struct {
  username Username
  width uint64
  center uint64
}
