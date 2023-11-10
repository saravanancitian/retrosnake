package main

import (
	"fmt"
	"retrosnake/retrosnake"

	"github.com/hajimehoshi/ebiten/v2"
)

const SCREEN_WIDTH int = 320
const SCREEN_HEIGHT int = 480

func main() {
	fmt.Print("Hello World")
	ebiten.SetWindowTitle("Retro Snake")
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.RunGame(retrosnake.NewApp())
}
