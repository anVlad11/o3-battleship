package game

import "errors"

func (game *Game) Shot(coordinate string) (ShotResult, error) {
	x, y, err := game.parseCoordinate(coordinate)
	if err != nil {
		return ShotResult{}, err
	}

	return game.shot(x, y)
}

func (game *Game) shot(x int64, y int64) (ShotResult, error) {
	shotResult := ShotResult{}

	if game.isEnded {
		return shotResult, errors.New("game is ended")
	}

	if int64(len(game.field)) < x || int64(len(game.field)) < y {
		return shotResult, errors.New("shot is outside of the game field")
	}

	game.gameChangesMutex.Lock()
	defer game.gameChangesMutex.Unlock()

	game.field[x][y].Shot = true

	if game.field[x][y].Ship != nil {
		shotResult.Knock = true

		shotResult.Destroy = true
		for _, dot := range game.field[x][y].Ship.Coordinates {
			if !dot.Shot {
				shotResult.Destroy = false
				break
			}
		}
	}

	if shotResult.Destroy == true {
		game.field[x][y].Ship.IsAlive = false
	}

	shotResult.End = true
	for _, ship := range game.ships {
		for _, dot := range ship.Coordinates {
			if !dot.Shot {
				shotResult.End = false
				break
			}
		}
		if !shotResult.End {
			break
		}
	}
	if shotResult.End {
		game.isEnded = true
	}

	game.shotCount++

	return shotResult, nil
}

type ShotResult struct {
	Destroy bool `json:"destroy"`
	Knock   bool `json:"knock"`
	End     bool `json:"end"`
}
