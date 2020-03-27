package main

import (
	"github.com/anvlad11/o3-battleship/game"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
)

var GameInstance game.Game

func main() {
	if os.Getenv("PREINIT") == "1" {
		initGame()
	}

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/create-matrix", CreateMatrix)
	r.POST("/ship", CreateShips)
	r.POST("/shot", CreateShot)
	r.POST("/clear", ClearGame)
	r.GET("/state", GetState)

	err := r.Run(os.Getenv("HTTP_HOST"))
	if err != nil {
		panic(err)
	}
}

func initGame() {
	GameInstance = game.NewGame(10)
	GameInstance.CreateShips("1A 2B,3D 3E")
	GameInstance.Shot("1A")
	GameInstance.Shot("1B")
	GameInstance.Shot("2A")
	GameInstance.Shot("2B")
	GameInstance.Shot("3A")
	GameInstance.GetStats()
}

type CreateMatrixRequest struct {
	Range int64 `json:"range"`
}

func CreateMatrix(ctx *gin.Context) {
	var err error
	request := CreateMatrixRequest{}

	err = ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(400, "Request body error")
		return
	}

	err = game.ValidateFieldSize(request.Range)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	GameInstance = game.NewGame(request.Range)

	ctx.JSON(200, nil)
}

type CreateShipsRequest struct {
	Coordinates string `json:"Coordinates"`
}

func CreateShips(ctx *gin.Context) {
	var err error
	request := CreateShipsRequest{}

	err = ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(400, "Request body error")
		return
	}

	err = GameInstance.CreateShips(request.Coordinates)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(200, nil)
}

func ClearGame(ctx *gin.Context) {
	GameInstance.Clear()

	ctx.JSON(200, nil)
}

func GetState(ctx *gin.Context) {
	stats := GameInstance.GetStats()

	ctx.JSON(200, stats)
}

type CreateShotRequest struct {
	Coord string `json:"coord"`
}

func CreateShot(ctx *gin.Context) {
	var err error
	request := CreateShotRequest{}

	err = ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(400, "Request body error")
		return
	}

	err = GameInstance.ValidateCoordinate(request.Coord)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	shotResult, err := GameInstance.Shot(request.Coord)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(200, shotResult)
}
