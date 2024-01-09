package structures

import "go.mongodb.org/mongo-driver/bson/primitive"

type BetData struct {
	PlayerID   string `json:"playerID"`
	RoomID     string `json:"roomID"`
	BetAmount1 int    `json:"betAmount1"`
	BetAmount2 int    `json:"betAmount2"`
}

type PlayerData struct {
	PlayerID            string  `json:"playerID"`
	BetAmount1          float64 `json:"betAmount1"`
	BetAmount2          float64 `json:"betAmount2"`
	CashoutAmount1      float64 `json:"cashoutAmount1"`
	CashoutAmount2      float64 `json:"cashoutAmount2"`
	HearbeatTimestamp   float64 `json:"hearbeatTimestamp"`
	IsBet1AmountCashout bool    `json:"isBet1AmountCashout"`
	IsBet2AmountCashout bool    `json:"isBet2AmountCashout"`
	IsPlayerConnected   bool    `json:"isPlayerConnected"`
	Multiplier          float64 `json:"multiplier"`
	PlayerImageID       string  `json:"playerImageID"`
	PlayerName          string  `json:"playerName"`
	WinAmount1          float64 `json:"winAmount1"`
	WinAmount2          float64 `json:"winAmount2"`
}

type RoomDetails struct {
	PlayerInRoom        []PlayerData `json:"PlayerInRoom"`
	WaitingPlayerInRoom []string     `json:"WaitingPlayerInRoom"`
	GameStatus          string       `json:"gameStatus"`
	IsGameStarted       bool         `json:"isGameStarted"`
	RoomID              string       `json:"roomID"`
	RoomStatus          string       `json:"roomStatus"`
}

type CashOutBetData struct {
	UserID  string  `json:"user_id"`
	RoomID  string  `json:"roomID"`
	BetNum  string  `json:"bet_num"`
	Cashout float64 `json:"cashout"`
}

type Player struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PlayerID     string             `bson:"playerID" json:"playerID"`
	WalletAmount float64            `bson:"walletAmount" json:"walletAmount"`
	AuthToken    string             `bson:"authToken" json:"authToken"`
}
