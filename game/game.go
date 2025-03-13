package main

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
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
    state.market = market

    if view == RegisterView {
      state.view = "wait"
    } else {
      state.room = rooms[user.room]

      if view == MakeView {
        state.view = "make"
      } else if view == CenterView {
        if rooms[user.room].username == username {
          state.view = "center"
        } else {
          state.view = "wait"
        }
      } else {
        state.view = "trade"
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

      if user.side == Long {
        profit = float64(round.TrueValue) - (float64(room.center) - float64(room.width) / 2)
      } else if user.side == Short {
        profit = (float64(room.width) / 2 + float64(room.center)) - float64(round.TrueValue)
      }

      stats := ranks[username]
      stats.pnl += profit
      ranks[username] = stats

      stats = ranks[room.username]
      stats.pnl -= profit
      ranks[room.username] = stats
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
      user.totalPlace += uint(curRank)
      user.totalPnl += stats.pnl
      users[stats.username] = user

      lastProfit = stats.pnl
      lastRoom = int(stats.room)
    }
  }

  fmt.Println(users)
}
