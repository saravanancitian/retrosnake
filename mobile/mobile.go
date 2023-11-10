package mobile

import (
	"retrosnake"

	"github.com/hajimehoshi/ebiten/v2/mobile"
)

////go run github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile  bind -target android -javapkg com.tictactoe.tictactoe -o ./bin/tictactoe.aar .\mobile

type IGameCallback interface {
	GameOverCallBack(winner int, duration int64)
}

var game *retrosnake.App

func init() {
	// yourgame.Game must implement ebiten.Game interface.
	// For more details, see
	// * https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Game

	game = retrosnake.NewApp()

	mobile.SetGame(game)
}

func RegisterGameCallback(callback IGameCallback) {
	game.RegisterIGameCallback(func(winner int, duration int64) { callback.GameOverCallBack(winner, duration) })
}

func PlayAgain(ngameplayed, nwin int) {
	game.PlayAgain(ngameplayed, nwin)
}

func Pause() {
	game.Pause()
}

func Resume() {
	game.Resume()
}

func Destroy() {
	game.Destroy()
}

func SetSoundOff(off bool) {
	game.SetSoundOff(off)
}

func SetShowTimerOff(off bool) {
	game.SetShowTimerOff(off)
}
