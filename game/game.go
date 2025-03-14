package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
)

var (
	users  = make(map[Username]User)
	view   = RegisterView
	rooms  []Room
	market string
	mtx    = sync.Mutex{}
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
    state.Side = user.side

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
					if user.side == None {
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

// requires mutex
func MakeGroups(NumRooms uint) []Room {
  arr := make([]Username, 0) 
  for user := range users {
    arr = append(arr, user)
  }

  rand.Shuffle(len(arr), func(i, j int) {
    arr[i], arr[j] = arr[j], arr[i]
  })

  rooms := make([]Room, NumRooms)

  for i, username := range arr {
    user := users[username] 
    user.room = uint(i) % NumRooms

    user.currPlace = 0
    user.currPnl = 0

    users[username] = user

    rooms[user.room].Ranks = append(rooms[user.room].Ranks, Rank{username, 0})
  }

  return rooms
}

func RunGame(rounds []Round) {
	mtx.Lock()
	defer mtx.Unlock()

	view = RegisterView
	mtx.Unlock()
	fmt.Println("press enter to begin game")
	waitForEnter()
	mtx.Lock()

  var rooms []Room

	// assign rooms
	NumRooms := uint(math.Ceil(float64(len(users)) / 10.0))

	for roundNum, round := range rounds {

    if roundNum % RoundsPerGroup == 0 {
      rooms = MakeGroups(NumRooms)
    }

		view = MakeView
		market = round.Market

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
			user.side = Long
			users[username] = user
			PushUserState(username)
		}

		mtx.Unlock()
		//time.Sleep(time.Second * 10)
		waitForEnter()
		mtx.Lock()

		// ranking
		for username, user := range users {
			var profit float64
			room := rooms[user.room]
			if username == room.Username || room.Username == "" {
				continue
			}

			if user.side == Long {
				profit = float64(round.TrueValue) - (float64(room.Center) + float64(room.Width)/2)
				fmt.Println("profit from Long: ", profit)
				fmt.Println("round.TrueValue: ", float64(round.TrueValue))
				fmt.Println("room.Center: ", float64(room.Center))
				fmt.Println("room.Width: ", float64(room.Width))
			} else if user.side == Short {
				profit = float64(round.TrueValue) + (float64(room.Width)/2 - float64(room.Center))
				fmt.Println("profit from Short: ", profit)
				fmt.Println("round.TrueValue: ", float64(round.TrueValue))
				fmt.Println("room.Center: ", float64(room.Center))
				fmt.Println("room.Width: ", float64(room.Width))
			}

      adj_profit := profit / math.Sqrt(round.TrueValue * room.Width)

      user.currPnl += adj_profit
      users[username] = user

      maker := users[room.Username]
      maker.currPnl -= adj_profit
      users[room.Username] = maker
		}

    for r := range rooms {
      room := &rooms[r]
      sort.Slice(room.Ranks, func(i, j int) bool {
        return users[room.Ranks[i].Username].currPnl > users[room.Ranks[j].Username].currPnl
      })

      next_place := uint(0)
      last_profit := math.Inf(-1)

      for i, stats := range room.Ranks {
        if users[stats.Username].currPnl != last_profit {
          next_place++
        }
        room.Ranks[i].Rank = next_place
        last_profit = users[stats.Username].currPnl

        user := users[stats.Username]
        user.currPlace = next_place
      }
    }

    if (roundNum + 1) % RoundsPerGroup == 0 { // total it
      for username, user := range users {
        user.totalPnl += user.currPnl
        user.totalPlace += user.currPlace
        users[username] = user
      }
    }

	}

	type Pair struct {
		Username Username
		User     User
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
		fmt.Println(i+1, ". ", pairs[i].Username)
	}
}
