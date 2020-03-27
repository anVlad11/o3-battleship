package game

import "fmt"

func (game *Game) GetStats() Stats {
	gameStats := Stats{
		ShipCount: int64(len(game.ships)),
		ShotCount: game.shotCount,
	}

	for _, ship := range game.ships {
		if !ship.IsAlive {
			gameStats.Destroyed++
		} else {
			for _, coordinate := range ship.Coordinates {
				if coordinate.Shot {
					gameStats.Knocked++
					break
				}
			}
		}
	}

	game.PrintStats()

	return gameStats
}

type Stats struct {
	ShipCount int64 `json:"ship_count"`
	Destroyed int64 `json:"destroyed"`
	Knocked   int64 `json:"knocked"`
	ShotCount int64 `json:"shot_count"`
}

//Рисует стейт игры в консоль, дебажная функция
func (game *Game) PrintStats() {
	fmt.Print("\t")
	for i := range game.field {
		fmt.Print(string(game.availableRunes[i]) + "\t")
	}
	fmt.Println()

	for i, slots := range game.field {
		fmt.Print(i + 1)
		fmt.Print("\t")
		for _, slot := range slots {
			if slot.Ship != nil {
				if slot.Ship.IsAlive == false {
					fmt.Print("X")
				} else if slot.Shot {
					fmt.Print("x")

				} else {
					fmt.Print("O")
				}
			} else if slot.Shot {
				fmt.Print(".")
			} else {
				fmt.Print(" ")
			}
			fmt.Print("\t")
		}
		fmt.Println()
	}
	fmt.Println()
}
