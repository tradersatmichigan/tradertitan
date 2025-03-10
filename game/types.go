package game

type Username = string

type RegisterResult = int

const (
  Success RegisterResult = iota
  UsernameTaken
)

type Side = int

const (
  Long Side = iota
  Short
  None
)

type User struct {
  Sse chan GameState 
  Room uint

  TotalPlacement uint
  TotalPnl int

  Side Side
}

type Quote struct {
  Username Username
  Price uint
}

type Trade struct {
  Quote Quote
  Side Side
}

type Room struct {
  Bid Quote
  Ask Quote
}

type GameState struct {
  Room Room
  Side Side
  Trading bool
}

type Round struct {
  True uint
  Market string
}
