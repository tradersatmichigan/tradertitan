package game

type registerArgs struct {
  username string
}

type registerResult = int

const (
  success registerResult = iota
  usernameTaken
)

type side = int

const (
  buy side = iota
  sell
)

type quote struct {
  username string
  price uint
}

type tradeArgs struct {
  quote quote
  side side
}

type sseArgs struct {
  sseSend chan <- gameState 
  username string
}

type gameState struct {
  bid quote 
  ask quote
  market string
  tradingOpen bool
}
