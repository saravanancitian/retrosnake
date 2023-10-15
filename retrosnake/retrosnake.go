package retrosnake

import "github.com/hajimehoshi/ebiten/v2"

type RetroSnake struct {
	GameOverCallBack func(int, int64)
}

func (r *RetroSnake) Update(delta int64) {

}

func (r *RetroSnake) Draw(screen *ebiten.Image) {

}

func NewRetroSnake(rm *ResourceManager, screenwidth int, screenheight int, callback func(int, int64)) *RetroSnake {
	var retrosnake = new(RetroSnake)
	retrosnake.Init(rm, screenwidth, screenheight, callback)
	return retrosnake
}

func (r *RetroSnake) Init(rm *ResourceManager, screenWidth int, screenHeight int, callback func(int, int64)) {

}
