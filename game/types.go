package main

type Username = string

type User struct {
	room       uint
	totalPnl   float64
	totalPlace uint
	currPlace  uint

	side     Side
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
	View   string  `json:"view"`
	Room   Room    `json:"room"`
	Market string  `json:"market"`
	Pnl    float64 `json:"pnl"`
	Place  uint    `json:"place"`
}

type Round struct {
	Market    string
	TrueValue int
}

type Room struct {
	Username Username `json:"username"`
	Width    uint64   `json:"width"`
	Center   uint64   `json:"center"`
}
