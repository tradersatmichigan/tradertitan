package game

import "math"


var (
  Advance chan any

  Register chan Username
  RegisterOut chan RegisterResult

  OpenTrade chan Trade
  HitTrade chan Trade
)

func init() {
  Advance = make(chan any)
  Register = make(chan Username)
  RegisterOut = make(chan RegisterResult)
  OpenTrade = make(chan Trade)
  HitTrade = make(chan Trade)
}

func RunGame(Rounds []Round) {
  users := make(map[Username] User)

  RegisterLoop:
  for {
    select {
    case <- Advance:
      break RegisterLoop
    case username := <- Register:
      if _, ok := users[username]; ok {
        RegisterOut <- UsernameTaken
      } else {
        users[username] = User{Sse : make(chan GameState, 1)}
        RegisterOut <- Success
      }
    }
  }

  // assign rooms
  numRooms := uint(math.Ceil(float64(len(users) / 10.0)))
  currentRoom := uint(0)
  for username, user := range users {
    user.Room = currentRoom % numRooms
    users[username] = user
    currentRoom++
  }

  // main game loop
  for round := range Rounds {

    rooms := make([]Room, numRooms)
    // reset state
    for username, user := range users {
      user.Side = None
      users[username] = user
    }
    
    OpenLoop:
    for {
      select {
      case <-Advance:
        break OpenLoop
      case <- HitTrade: // discard to prevent accidental trade on next round
      case trade := <- OpenTrade:
        if user, ok := users[trade.Quote.Username]; ok {
          if (trade.Quote.Price > rooms[user.Room].Bid.Price || rooms[user.Room].Bid.Username == "") && 
             (trade.Quote.Price < rooms[user.Room].Ask.Price || rooms[user.Room].Ask.Username == "") {
            if trade.Side == Long {
              rooms[user.Room].Bid.Username = trade.Quote.Username
              rooms[user.Room].Bid.Price = trade.Quote.Price
            } else if trade.Side == Short {
              rooms[user.Room].Ask.Username = trade.Quote.Username
              rooms[user.Room].Ask.Price = trade.Quote.Price
            }
            // broadcast change
          }
        }
      }
    }

    HitLoop:
    for {
      select {
      case <- Advance:
        break HitLoop
      case <- OpenTrade:
      case trade := <- HitTrade:
        if user, ok := users[trade.Quote.Username]; ok {
          user.Side = trade.Side
          users[trade.Quote.Username] = user
        }
      }
    }

    // total up pnl and rank
  }
}

