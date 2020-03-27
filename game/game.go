package game

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

const runes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Game struct {
	field          [][]*Slot
	isStarted      bool
	isEnded        bool
	ships          []*Ship
	availableRunes string

	gameChangesMutex sync.Mutex

	shotCount int64
}

type Slot struct {
	Ship *Ship
	Shot bool
}

func NewGame(gameSize int64) Game {
	game := Game{}
	game.Init(gameSize)

	return game
}

var initMutex sync.Mutex

func (game *Game) Init(gameSize int64) {
	initMutex.Lock()
	defer initMutex.Unlock()

	field := make([][]*Slot, gameSize)
	runes := make([]rune, gameSize)

	for i := range field {
		field[i] = make([]*Slot, gameSize)
		for j := range field[i] {
			field[i][j] = &Slot{}
		}

		runes[i] = toChar(i + 1)
	}

	game.field = field
	game.isStarted = false
	game.isEnded = false
	game.ships = []*Ship{}
	game.availableRunes = string(runes)
	game.gameChangesMutex = sync.Mutex{}
	game.shotCount = 0

}

func (game *Game) Clear() {
	game.Init(int64(len(game.field)))
}

func ValidateFieldSize(fieldSize int64) error {
	if fieldSize < 0 {
		return errors.New("field size is too small")
	}

	if fieldSize > int64(len(runes)) {
		return errors.New("field size is too large, maximum is " + strconv.FormatInt(int64(len(runes)), 10))
	}

	return nil
}

func (game *Game) ValidateCoordinate(coordinate string) error {
	_, _, err := game.parseCoordinate(coordinate)

	return err
}

func (game *Game) parseCoordinate(coordinate string) (int64, int64, error) {
	coordinateLength := len(coordinate)
	if coordinateLength < 2 {
		return 0, 0, errors.New("coordinate count is invalid")
	}

	yChar := string(coordinate[coordinateLength-1])
	y := int64(strings.Index(game.availableRunes, yChar))
	if y == -1 {
		return 0, 0, errors.New("symbolical coordinate is invalid")
	}

	xChars := coordinate[:coordinateLength-1]

	x, err := strconv.ParseInt(xChars, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	//visual coordinate x = stored coordinate x + 1
	return x - 1, y, nil
}

func toChar(i int) rune {
	return rune('A' - 1 + i)
}
