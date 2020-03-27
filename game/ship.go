package game

import (
	"errors"
	"strings"
)

type Ship struct {
	Coordinates []*Slot
	IsAlive     bool
}

func (game *Game) createShip(coordinates string) error {
	pairs := strings.Split(coordinates, " ")
	if len(pairs) < 1 || len(pairs) > 2 {
		return errors.New("ship coordinates are invalid")
	}

	top, left, err := game.parseCoordinate(pairs[0])
	if err != nil {
		return err
	}

	bottom, right, err := game.parseCoordinate(pairs[1])
	if err != nil {
		return err
	}

	err = game.validateShipCoordinates(top, left, bottom, right)
	if err != nil {
		return err
	}

	ship := &Ship{
		Coordinates: []*Slot{},
		IsAlive:     true,
	}

	game.placeShip(ship, top, left, bottom, right)
	game.ships = append(game.ships, ship)
	game.isStarted = true

	return nil
}

func (game *Game) placeShip(ship *Ship, top int64, left int64, bottom int64, right int64) {
	for i := top; i <= bottom; i++ {
		for j := left; j <= right; j++ {
			game.field[i][j].Ship = ship
			ship.Coordinates = append(ship.Coordinates, game.field[i][j])
		}
	}
}

func (game *Game) validateShipCoordinates(top int64, left int64, bottom int64, right int64) error {
	if top > bottom || left > right {
		return errors.New("invalid ship form")
	}

	if top < 0 || left < 0 || bottom >= int64(len(game.field)) || right >= int64(len(game.field)) {
		return errors.New("coordinates out of game field")
	}

	for i := top; i <= bottom; i++ {
		for j := left; j <= right; j++ {
			if game.field[i][j].Ship != nil {
				return errors.New("place already occupied")
			}
		}
	}

	return nil
}

func (game *Game) CreateShips(coordinates string) error {
	if game.isStarted {
		return errors.New("game already started")
	}

	game.gameChangesMutex.Lock()
	defer game.gameChangesMutex.Unlock()

	coordinatePairs := strings.Split(coordinates, ",")
	for _, pair := range coordinatePairs {
		err := game.createShip(pair)
		if err != nil {
			game.Clear()
			return err
		}
	}

	return nil
}
