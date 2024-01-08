package controller

import (
	"aviator/collections"
	configs "aviator/config"
	"aviator/constant"
	structure "aviator/structures"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	// "github.com/labstack/gommon/bytes"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

var requestQue = []structure.BetData{}
var cashOutQue = []structure.CashOutBetData{}
var thresholdValue float64

type Headers struct {
	Authorization string `json:"Authorization"`
}
type BetInfo struct {
	UserID    string  `json:"user_id"`
	BetAmount float64 `json:"bet_amount"`
	BetType   string  `json:"bet_type"`
}

type AviatorWinResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Rewards string `json:"rewards"`
		Winning string `json:"winning"`
	} `json:"data"`
}

type reduceResult struct {
	totalBetAmount     float64
	maxAmount          float64
	totalCashoutAmount float64
}

var response = map[string]any{
	"status":          true,
	"data":            []string{},
	"threshold_value": thresholdValue,
	"message":         "Cash Out Successfully",
}

/* <----------------- PLACE BET ---------------------> */
func PlaceBet(c echo.Context) error {
	var requestBody structure.BetData
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}
	requestQue = append(requestQue, requestBody)
	processBetQue(&requestQue)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"data":    []string{},
		"message": "Place Bet Successfully",
	})
}

func processBetQue(que *[]structure.BetData) {
	betQue := *que
	for _, v := range betQue {
		processBetRequest(v)
	}
}

func processBetRequest(betData structure.BetData) {
	key := fmt.Sprintf("%s:%s", constant.AVIATOR_ROOM, betData.RoomID)
	result, err := configs.GetString(key)
	if err != nil {
		log.Println("Error At process")
		return
	}
	var roomDetails structure.RoomDetails
	err = json.Unmarshal([]byte(result), &roomDetails)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", constant.RED(err.Error()))
		return
	}
	playerData, err := findPlayerInRoom(roomDetails.PlayerInRoom, betData.PlayerID, false)
	if err != nil {
		fmt.Println("Error:", constant.RED(err.Error()))
		return
	}
	if betData.BetAmount1 > 0 {
		playerData.BetAmount1 = float64(betData.BetAmount1)
	}
	if betData.BetAmount2 > 0 {
		playerData.BetAmount2 = float64(betData.BetAmount2)
	}
	updatePlayerInRoom(roomDetails, playerData, key)
}

func findPlayerInRoom(people []structure.PlayerData, playerId string, cashout bool) (structure.PlayerData, error) {
	for _, p := range people {
		if p.PlayerID == playerId {
			if cashout {
				cashOutQue = cashOutQue[1:]
			} else {
				requestQue = requestQue[1:]
			}
			return p, nil
		}
	}
	return structure.PlayerData{}, fmt.Errorf("Player with ID %s not found", playerId)
}

func updatePlayerInRoom(roomDetails structure.RoomDetails, playerData structure.PlayerData, key string) {
	var newPlayers []structure.PlayerData
	for _, v := range roomDetails.PlayerInRoom {
		if v.PlayerID == playerData.PlayerID {
			newPlayers = append(newPlayers, playerData)
		} else {
			newPlayers = append(newPlayers, v)
		}
	}

	roomDetails.PlayerInRoom = newPlayers
	jsonData, err := json.Marshal(roomDetails)
	if err != nil {
		fmt.Println("Error:", constant.RED(err.Error()))
		return
	}
	jsonString := string(jsonData)
	if err := configs.SetString(key, jsonString); err != nil {
		fmt.Println("Error:", constant.RED(err.Error()))
	}
}

/* <----------------- CASH OUT ---------------------> */

func CashOut(c echo.Context) error {
	var requestBody structure.CashOutBetData
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}
	cashOutQue = append(cashOutQue, requestBody)
	processCashOutQue()
	return c.JSON(http.StatusOK, response)
}

func processCashOutQue() {
	for _, v := range cashOutQue {
		processCashOutRequest(v)
	}
}

func processCashOutRequest(playerObject structure.CashOutBetData) {
	key := fmt.Sprintf("%s:%s", constant.AVIATOR_ROOM, playerObject.RoomID)
	result, err := configs.GetString(key)
	var roomDetails structure.RoomDetails
	err = json.Unmarshal([]byte(result), &roomDetails)
	if err != nil {
		log.Println(constant.RED(err.Error()))
		return
	}
	playerData, err := findPlayerInRoom(roomDetails.PlayerInRoom, playerObject.UserID, true)
	if err != nil {
		log.Println(constant.RED(err.Error()))
		return
	}
	isUserCashout := threshold(playerObject.RoomID, playerObject.Cashout, roomDetails, playerData, playerObject, key)
	if isUserCashout && playerObject.BetNum == "bet1" && playerObject.Cashout != 0 {
		updateUserWalletAmount(playerObject)
	}
	if isUserCashout && playerObject.BetNum == "bet2" && playerObject.Cashout != 0 {
		updateUserWalletAmount(playerObject)
	}
}

func threshold(roomId string, cashout float64, roomDetails structure.RoomDetails, playerData structure.PlayerData, playerObject structure.CashOutBetData, key string) bool {
	if len(roomDetails.PlayerInRoom) > 0 {
		if playerObject.BetNum == "bet1" {
			playerData.IsBet1AmountCashout = true
			playerData.CashoutAmount1 = cashout
			playerData.WinAmount1 = cashout
			updatePlayerInRoom(roomDetails, playerData, key)
		}
		if playerObject.BetNum == "bet2" {
			playerData.IsBet2AmountCashout = true
			playerData.CashoutAmount2 = cashout
			playerData.WinAmount2 = cashout
			updatePlayerInRoom(roomDetails, playerData, key)
		}
		result, err := configs.GetString(key)
		var roomDetail structure.RoomDetails
		err = json.Unmarshal([]byte(result), &roomDetail)
		if err != nil {
			log.Println(constant.RED(err.Error()))
			return false
		}
		res := playersInRoomArrayReduce(roomDetail.PlayerInRoom)
		var distributedAmount = 0.8 * res.totalBetAmount
		thresholdValue = (distributedAmount - res.totalCashoutAmount) / res.maxAmount
		if distributedAmount < res.totalCashoutAmount {
			response["status"] = false
			return false
		} else {
			return true
		}
	}
	return false
}

func playersInRoomArrayReduce(playersInRoomArray []structure.PlayerData) reduceResult {
	// Initial accumulator values
	acc := reduceResult{}

	// Reduce function
	for _, obj := range playersInRoomArray {
		acc.totalBetAmount += obj.BetAmount1 + obj.BetAmount2
		acc.maxAmount = max(
			acc.maxAmount,
			obj.BetAmount1,
			obj.BetAmount2,
		)
		acc.totalCashoutAmount += cashoutAmount(obj.IsBet1AmountCashout, obj.CashoutAmount1) +
			cashoutAmount(obj.IsBet2AmountCashout, obj.CashoutAmount2)
	}

	return acc
}

func max(values ...float64) float64 {
	result := values[0]
	for _, v := range values[1:] {
		if v > result {
			result = v
		}
	}
	return result
}

func cashoutAmount(isCashout bool, amount float64) float64 {
	if isCashout {
		return amount
	}
	return 0
}

func postWalletAmount(apiUrl string, postData BetInfo, token string) (AviatorWinResponse, error) {

	method := "POST"

	payload, err := json.Marshal(postData)
	if err != nil {
		return AviatorWinResponse{}, err
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, apiUrl, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println(err)
		return AviatorWinResponse{}, err
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return AviatorWinResponse{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return AviatorWinResponse{}, err
	}
	var response AviatorWinResponse
	err = json.Unmarshal([]byte(string(body)), &response)
	if err != nil {
		fmt.Println(err)
		return AviatorWinResponse{}, err
	}
	return response, nil
}

func updateUserWalletAmount(playerObject structure.CashOutBetData) {
	apiUrl := os.Getenv("POST_API_WALLET_URL")

	betInfo := BetInfo{
		UserID:    playerObject.UserID,
		BetAmount: playerObject.Cashout,
		BetType:   "win",
	}
	filter := bson.D{{Key: "playerID", Value: playerObject.UserID}}
	var player structure.Player
	err := collections.PLAYER_WALLET_AMOUNT.FindOne(context.TODO(), filter).Decode(&player)
	if err != nil {
		return
	}
	var authToken = player.AuthToken
	postWalletData, err := postWalletAmount(apiUrl, betInfo, authToken)
	if err != nil {
		fmt.Println(err)
		return
	}
	if postWalletData.Data.Rewards != "" || postWalletData.Data.Winning != "" {
		winning, err := strconv.ParseFloat(postWalletData.Data.Winning, 64)
		if err != nil {
			fmt.Println("Error parsing winning:", err)
			return
		}
		rewards, err := strconv.ParseFloat(postWalletData.Data.Rewards, 64)
		if err != nil {
			fmt.Println("Error parsing rewards:", err)
			return
		}
		update := bson.M{"$set": bson.M{"walletAmount": rewards + winning}}
		_, err = collections.PLAYER_WALLET_AMOUNT.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Fatal(err)
		}
	}
}
