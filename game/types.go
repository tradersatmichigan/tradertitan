package main

type Username = string

type User struct {
  room uint
  totalPnl float64
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
)

type Side = int

const (
  Long Side = iota
  Short
  None
)

type GameState struct {
  view string
  room Room
  market string
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
