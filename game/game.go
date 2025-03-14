package main

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

var (
  users = make(map[Username] User)
  view = RegisterView
  rooms []Room
  market string
  mtx = sync.Mutex{}
)
func Register(username Username) bool {
  mtx.Lock()
  defer mtx.Unlock()

  if view != RegisterView {
    return false
  } else if _, ok := users[username]; ok {
    return false
  } 

  users[username] = User{}

  return true
}

func GetUserChan(username Username) chan GameState {
  mtx.Lock()
  defer mtx.Unlock()

  if user, ok := users[username]; ok {
    user.datachan = make(chan GameState, 1)
    users[username] = user
    PushUserState(username)
    return user.datachan
  } else {
    return nil
  }
}

// require mutex
func PushUserState(username Username) {

  if user, ok := users[username]; ok {
    state := GameState{}
    state.Market = market
    state.Pnl = user.totalPnl
    state.Place = user.currPlace

    if view == RegisterView {
      state.View = "wait"
    } else {
      state.Room = rooms[user.room]

      if view == MakeView {
        state.View = "make"
      } else if view == CenterView {
        if rooms[user.room].Username == username {
          // state.Room.Center = 0
          state.View = "center"
        } else {
          state.View = "wait"
        }
      } else if view == TradeView {
        if rooms[user.room].Username == username {
          state.View = "wait"
        } else {
          if (user.side == None) {
            user.side = Long
            users[username] = user
            state.View = "trade"
          }
        }
      }
    }

    select {
    case user.datachan <- state:
    default:
    }
  }
}

func RunGame(rounds []Round) {
  mtx.Lock()
  defer mtx.Unlock()

  view = RegisterView
  mtx.Unlock()
  fmt.Println("press enter to begin game")
  waitForEnter()
  mtx.Lock()

  // assign rooms
  NumRooms := uint(math.Ceil(float64(len(users)) / 10.0))
  i := uint(0)
  for username, user := range users {
    user.room = i % NumRooms
    i++
    users[username] = user
  }

  for _, round := range rounds {
    view = MakeView
    market = round.Market
    rooms = make([]Room, NumRooms)

    for user := range users {
      PushUserState(user)
    }

    mtx.Unlock()
    //time.Sleep(time.Second * 10)
    waitForEnter()
    mtx.Lock()

    view = CenterView
    for user := range users {
      PushUserState(user)
    }

    mtx.Unlock()
    //time.Sleep(time.Second * 10)
    waitForEnter()
    mtx.Lock()

    view = TradeView
    for username, user := range users {
      user.side = None
      users[username] = user
      PushUserState(username)
    }

    mtx.Unlock()
    //time.Sleep(time.Second * 10)
    waitForEnter()
    mtx.Lock()

    // ranking
    type UserRanker struct {
      username Username
      pnl float64
      room uint
    }

    ranks := make(map[Username] UserRanker)
    for username, user := range users {
      ranks[username] = UserRanker{username, 0, user.room}
    }

    for username, user := range users {
      var profit float64
      room := rooms[user.room]
      if (username == room.Username) {
        continue;
      }

      if user.side == Long {
        profit = float64(round.TrueValue) - (float64(room.Center) + float64(room.Width) / 2)
        fmt.Println("profit from Long: ", profit)
        fmt.Println("round.TrueValue: ", float64(round.TrueValue))
        fmt.Println("room.Center: ", float64(room.Center))
        fmt.Println("room.Width: ", float64(room.Width))
      } else if user.side == Short {
        profit = float64(round.TrueValue) + (float64(room.Width) / 2 - float64(room.Center))
        fmt.Println("profit from Short: ", profit)
        fmt.Println("round.TrueValue: ", float64(round.TrueValue))
        fmt.Println("room.Center: ", float64(room.Center))
        fmt.Println("room.Width: ", float64(room.Width))
      }

      stats := ranks[username]
      stats.pnl += profit
      ranks[username] = stats

      stats = ranks[room.Username]
      stats.pnl -= profit
      ranks[room.Username] = stats
    }

    sorter := make([]UserRanker, 0)

    for _, user := range ranks {
      sorter = append(sorter, user)
    }

    sort.Slice(sorter, func(i, j int) bool {
      if sorter[i].room != sorter[j].room {
        return sorter[i].room < sorter[j].room
      }
      return sorter[i].pnl > sorter[j].pnl
    })

    lastRoom := -1
    lastProfit := math.Inf(-1)
    curRank := 0
    for _, stats := range sorter {
      if stats.room != uint(lastRoom) {
        curRank = 0
      } else if stats.pnl != lastProfit {
        curRank++
      }

      user := users[stats.username]
      user.currPlace = uint(curRank + 1)
      user.totalPlace += uint(curRank)
      user.totalPnl += stats.pnl
      users[stats.username] = user

      lastProfit = stats.pnl
      lastRoom = int(stats.room)
    }
  }

  type Pair struct {
    Username Username
    User User
  }
  
  pairs := make([]Pair, 0)
  for username, user := range users {
    pairs = append(pairs, Pair{username, user})
  }

  sort.Slice(pairs, func(i, j int) bool {
    if pairs[i].User.totalPlace != pairs[j].User.totalPlace {
      return pairs[i].User.totalPlace > pairs[j].User.totalPlace
    }
    return pairs[i].User.totalPnl > pairs[j].User.totalPnl
  })

  for i := range pairs {
    fmt.Println(i + 1, ". ", pairs[i].Username)
  }
}
