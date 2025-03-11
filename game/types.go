package main

type Username = string

type RegisterResult = int
const (
  Success RegisterResult = iota
  UsernameTaken
)

type User struct {
  Room uint
  TotalPlacement uint
  TotalPnl int

  Side Side
  CurPnl int
}

type Round struct {
  Market string
  TrueVal int
}

type Quote struct {
  Username Username
  Center   int
  Width    uint
}

type Side = int
const (
  Long Side = iota
  Short
  None
)

type SpreadArgs struct {
  Username Username
  Width uint
}

type CenterArgs struct {
  Username Username
  Center int
}

type TradeArgs struct {
  Username Username
  Side Side
}

type Display struct {
  view string
  market string
  Quote Quote
}
