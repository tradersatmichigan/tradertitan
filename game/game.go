package main

import (
	"fmt"
	"math"
	"sort"
)

const buffer = 10

var (
  Advance chan any
  
  Register chan Username
  RegisterReply chan RegisterResult

  Spread chan SpreadArgs
  Center chan CenterArgs
  Trade chan TradeArgs

  GetState chan Username
  DisplayOut chan Display
)

func init() {
  Advance = make(chan any)
  Register = make(chan Username)
  RegisterReply = make(chan RegisterResult)
  Spread = make(chan SpreadArgs)
  Center = make(chan CenterArgs)
  Trade = make(chan TradeArgs)
  GetState = make(chan Username)
  DisplayOut = make(chan Display)
}

func RunGame(Rounds []Round) {
  Users := make(map[Username] User)
  
  RegisterLoop:
  for {
    select {
    case <-Advance:
      break RegisterLoop
    case <-GetState:
      DisplayOut <- Display{view: "wait"}
    case username := <- Register:
      if _, ok := Users[username]; ok {
        RegisterReply <- UsernameTaken
      } else {
        Users[username] = User{}
        RegisterReply <- Success
      }
    }
  }

  // assign rooms
  NumRooms := uint(
    math.Ceil(float64(len(Users)) / 10.0),
  )

  i := uint(0)
  for username, user := range Users {
    user.Room = i % NumRooms
    Users[username] = user
    i++
  }

  // todo service gets
  for _, Round := range Rounds {
    
    Quotes := make([]Quote, NumRooms)

    MarketLoop:
    for {
      select {
      case <- Advance:
        break MarketLoop
      case user := <- GetState:
        DisplayOut <- Display{"make", Round.Market, Quotes[Users[user].Room]}
      case <- Center:
      case <- Trade:
      case args := <- Spread:
        if user, ok := Users[args.Username]; ok {
          quote := &Quotes[user.Room]
          if quote.Width > args.Width || quote.Username == "" {
            quote.Width = args.Width
            quote.Username = args.Username
          }
        }
      }
    }

    CenterLoop:
    for {
      select {
      case <-Advance:
        break CenterLoop
      case user := <- GetState:
        DisplayOut <- Display{"center", Round.Market, Quotes[Users[user].Room]}
      case <- Trade:
      case <- Spread:
      case args := <- Center:
        if user, ok := Users[args.Username]; ok {
          quote := &Quotes[user.Room]
          if quote.Username == args.Username {
            quote.Center = args.Center
          }
        }
      }
    }
    
    // reset trade state
    for name, user := range Users {
      user.Side = None
      user.CurPnl = 0
      Users[name] = user
    }

    TradeLoop:
    for {
      select {
      case <-Advance:
        break TradeLoop
      case user := <- GetState:
        DisplayOut <- Display{"trade", Round.Market, Quotes[Users[user].Room]}
      case <- Spread:
      case <- Center:
      case args := <- Trade:
        if user, ok := Users[args.Username]; ok && 
        Quotes[user.Room].Username != "" {
          user.Side = args.Side
          Users[args.Username] = user
        }
      }
    }

    // calculate round profit
    for name, user := range Users {
      quote := Quotes[user.Room]
      var profit int

      if user.Side == Long {
        profit = Round.TrueVal - (quote.Center + int(quote.Width))
      } else if user.Side == Short {
        profit = (quote.Center - int(quote.Width)) - Round.TrueVal
      }
      
      user.CurPnl += profit
      Users[name] = user

      op := Quotes[user.Room].Username
      opUser := Users[op]
      opUser.CurPnl -= profit
      Users[op] = opUser
    }

    // calculate room placement
    type RoomInfo struct {
      count uint
      last int
    }
    Rooms := make([]RoomInfo, NumRooms)
    for i := range Rooms {
      Rooms[i].last = math.MinInt
    }

    type userSorter struct {
      username Username
      profit int
      room uint
    }
    Arr := make([]userSorter, 0)
    for name, user := range Users {
      Arr = append(Arr, userSorter{name, user.CurPnl, user.Room})
    }
    sort.Slice(Arr, func(i, j int) bool {
      return Arr[i].profit < Arr[j].profit
    })

    for _, user := range Arr {
      if Rooms[user.room].last != user.profit {
        Rooms[user.room].last = user.profit
        Rooms[user.room].count++
      }
      
      oUser := Users[user.username]
      oUser.TotalPlacement += Rooms[user.room].count
      oUser.TotalPnl += oUser.CurPnl
      Users[user.username] = oUser
    }
  }

  fmt.Print(Users)
}
