package retrosnake

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"retrosnake/retrosnake/input"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	KEY_UP    int = 38
	KEY_DOWN  int = 40
	KEY_LEFT  int = 37
	KEY_RIGHT int = 39

	STATE_GAME_OVER_HALT int = 0
	STATE_INIT_NEW_GAME  int = 1
	STATE_GAME_RUNNING   int = 2
	STATE_GAME_PAUSED    int = 3
	STATE_GAME_OVER      int = 4

	LEFT_MARGIN int = 24
	TOP_MARGIN  int = 42

	GRID_BORDER int = 1

	GRID_WIDTH  int = 270
	GRID_HEIGHT int = 330

	CELL_WIDTH  int = 30
	CELL_HEIGHT int = 30

	ANIM_DELAY    int64 = 500
	SEC_IN_MILLIS int64 = 1000

	gameOverStr string = "GAME OVER"
	pausedStr   string = "PAUSED"
)

type RetroSnake struct {
	GameOverCallBack func(int, int64)

	rm *ResourceManager

	state int

	playAreaX      int
	playAreaY      int
	playAreaWidth  int
	playAreaHeight int

	head *SnakePart

	maxRow int
	maxCol int

	timeBg  *ebiten.Image
	bgImg   *ebiten.Image
	overlay *ebiten.Image
	about   *ebiten.Image

	upArrow     *ebiten.Image
	downArrow   *ebiten.Image
	leftArrow   *ebiten.Image
	rightArrow  *ebiten.Image
	yellowRound *ebiten.Image
	greenRound  *ebiten.Image

	headDown  *ebiten.Image
	headLeft  *ebiten.Image
	headRight *ebiten.Image
	headUp    *ebiten.Image

	mouseImg     *ebiten.Image
	resumeButton *ebiten.Image
	pauseButton  *ebiten.Image
	newButton    *ebiten.Image

	reverseDesl *ebiten.Image
	reverseSel  *ebiten.Image

	TOONEYNO_ttf *opentype.Font
	gameFont     font.Face

	random *rand.Rand

	mousePos image.Point

	totalMouseCapture int
	isPaused          bool
	check_game        bool
	isReverse         bool

	totalTime int64
	prevTime  int64

	buttonW int
	buttonH int
	bgW     int
	bgH     int

	timerBgX int
	timerBgY int
	timerBgW int
	timerBgH int

	timerStrX int
	timerStrY int

	txtTimer  string
	showTimer bool
	showMsg   bool

	upButtonX      int
	upButtonY      int
	downButtonX    int
	downButtonY    int
	leftButtonX    int
	leftButtonY    int
	rightButtonX   int
	rightButtonY   int
	reverseButtonX int
	reverseButtonY int
	reverseButtonW int
	reverseButtonH int
	newButtonX     int
	newButtonY     int
	pauseButtonX   int
	pauseButtonY   int
	npButtonW      int
	npButtonH      int
	controlPanelY  int

	gONewButtonX  int
	gONewButtonY  int
	resumeButtonX int
	resumeButtonY int

	playtime     int64
	playtimecalc int64

	animtime    int64
	isAnimating bool
	isGameover  bool

	settingTimerOff bool
	settingSoundOff bool
}

func (rs *RetroSnake) SetSoundOff(off bool) {
	rs.settingSoundOff = off
}

func (rs *RetroSnake) SetShowTimerOff(off bool) {
	rs.settingTimerOff = off
}

func (rs *RetroSnake) SetCallback(callback func(int, int64)) {
	rs.GameOverCallBack = callback
}

func (rs *RetroSnake) Init(rm *ResourceManager, screenWidth int, screenHeight int, callback func(int, int64)) {

	rs.rm = rm
	rs.GameOverCallBack = callback

	rs.LoadImages()
	rs.LoadFonts()

	rs.buttonW = 29
	rs.buttonH = 29
	rs.bgW = 320
	rs.bgH = 480
	rs.reverseButtonW = 40
	rs.reverseButtonH = 41

	rs.npButtonW = 78
	rs.npButtonH = 27

	rs.timerBgW = 93
	rs.timerBgH = 35

	rs.playAreaX = LEFT_MARGIN + GRID_BORDER
	rs.playAreaY = TOP_MARGIN + GRID_BORDER
	rs.playAreaWidth = GRID_WIDTH
	rs.playAreaHeight = GRID_HEIGHT

	rs.maxRow = (GRID_HEIGHT - 2) / CELL_HEIGHT
	rs.maxCol = (GRID_WIDTH - 2) / CELL_WIDTH

	rs.timerBgX = LEFT_MARGIN
	rs.timerBgY = rs.playAreaY - rs.timerBgH - 1

	//baseLine := rs.timerBgH - ((rs.timerBgH - 17) >> 1)

	rs.timerStrX = rs.timerBgX + 9
	rs.timerStrY = rs.timerBgY + 2 //baseLine - 2

	rs.txtTimer = "00:00"

	rs.controlPanelY = rs.playAreaY + rs.playAreaHeight + 2
	rs.upButtonY = rs.controlPanelY + 4

	rs.reverseButtonY = rs.upButtonY + rs.buttonH
	rs.rightButtonY = rs.reverseButtonY + ((rs.reverseButtonH - rs.buttonH) >> 1)
	rs.leftButtonY = rs.rightButtonY
	rs.downButtonY = rs.reverseButtonY + rs.reverseButtonH

	rs.rightButtonX = rs.bgW - (rs.buttonW << 1)
	rs.reverseButtonX = rs.rightButtonX - rs.reverseButtonW
	rs.downButtonX = rs.reverseButtonX + ((rs.reverseButtonW - rs.buttonW) >> 1)

	rs.upButtonX = rs.downButtonX

	rs.leftButtonX = rs.reverseButtonX - rs.buttonW

	rs.pauseButtonX = rs.leftButtonX - rs.npButtonW - (rs.buttonW >> 1)
	rs.pauseButtonY = rs.leftButtonY //controlPanelY + ((bgH - controlPanelY) >> 1);

	rs.newButtonX = rs.pauseButtonX - rs.npButtonW - (rs.buttonW >> 2)
	rs.newButtonY = rs.leftButtonY

	// application.scheduleTask(TASKID_INITGAME, 100, 0)

	rs.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	rs.state = STATE_INIT_NEW_GAME

}

func (rs *RetroSnake) SetDelay(delay int64) {
	rs.animtime = delay

}

func (rs *RetroSnake) DelayElapsed(delta int64) bool {
	if rs.animtime-delta > 0 {
		rs.animtime -= delta
	} else {
		return true
	}
	return false

}

func (rs *RetroSnake) initGame() {
	rs.isPaused = false
	rs.check_game = false
	rs.isReverse = false
	rs.totalMouseCapture = 0
	rs.totalTime = 0
	rs.prevTime = 0
	rs.txtTimer = "00:00"
	rs.CreateSnake()
	rs.GenerateMouse()
	// task_move_snake = application.scheduleTask(TASKID_MOVESNAKE, 500, 250)
	rs.StartTime()
	rs.state = STATE_GAME_RUNNING
	rs.SetDelay(ANIM_DELAY)

}

func (rs *RetroSnake) OnClick(x, y int) bool {
	var handled bool = false

	if (x > rs.upButtonX && x < (rs.upButtonX+rs.buttonW)) && (y > rs.upButtonY && y < (rs.upButtonY+rs.buttonH)) {
		handled = rs.OnKey(KEY_UP)
	} else if (x > rs.downButtonX && x < (rs.downButtonX+rs.buttonW)) && (y > rs.downButtonY && y < (rs.downButtonY+rs.buttonH)) {
		handled = rs.OnKey(KEY_DOWN)
	} else if (x > rs.leftButtonX && x < (rs.leftButtonX+rs.buttonW)) && (y > rs.leftButtonY && y < (rs.leftButtonY+rs.buttonH)) {
		handled = rs.OnKey(KEY_LEFT)
	} else if (x > rs.rightButtonX && x < (rs.rightButtonX+rs.buttonW)) && (y > rs.rightButtonY && y < (rs.rightButtonY+rs.buttonH)) {
		handled = rs.OnKey(KEY_RIGHT)
	} else if (x > rs.reverseButtonX && x < (rs.reverseButtonX+rs.reverseButtonW)) && (y > rs.reverseButtonY && y < (rs.reverseButtonY+rs.reverseButtonH)) {
		rs.isReverse = !rs.isReverse
		handled = true
	}

	return handled
}

func (rs *RetroSnake) OnKey(keyCode int) bool {

	var handled bool = false

	if rs.state == STATE_GAME_RUNNING {
		switch keyCode {
		case KEY_LEFT:
			var direction int
			if rs.isReverse {
				direction = DIRECTION_RIGHT
			} else {
				direction = DIRECTION_LEFT
			}
			handled = rs.head.changeDirection(direction)

		case KEY_RIGHT:

			var direction int
			if rs.isReverse {
				direction = DIRECTION_LEFT
			} else {
				direction = DIRECTION_RIGHT
			}
			handled = rs.head.changeDirection(direction)

		case KEY_UP:
			var direction int
			if rs.isReverse {
				direction = DIRECTION_DOWN
			} else {
				direction = DIRECTION_UP
			}
			handled = rs.head.changeDirection(direction)
		case KEY_DOWN:
			var direction int
			if rs.isReverse {
				direction = DIRECTION_UP
			} else {
				direction = DIRECTION_DOWN
			}
			handled = rs.head.changeDirection(direction)

		}
	}

	return handled
}

func (rs *RetroSnake) CalculatePlayTime(delta int64) {
	if !rs.isGameover {
		if rs.playtimecalc >= SEC_IN_MILLIS {
			rs.playtime += rs.playtimecalc
			rs.playtimecalc = delta

			tsec := rs.playtime / SEC_IN_MILLIS
			min := tsec / 60
			sec := tsec % 60
			rs.txtTimer = fmt.Sprintf("\n%02d:%02d", min, sec)

		} else {
			rs.playtimecalc += delta
		}
	}
}

func (rs *RetroSnake) StartTime() {
	// time_task =  application.scheduleTask(TASKID_TIME, 1000, 1000);
	// prevTime = System.currentTimeMillis();
}

func (rs *RetroSnake) StopTime() {
	// if(time_task != null)
	// 	time_task.cancel();
}

func (rs *RetroSnake) gamemove() {
	if rs.check_game {

		if rs.head.CheckTailBite() {
			rs.state = STATE_GAME_OVER
		}

		if rs.head.CheckCaptureMouse(rs.mousePos.Y, rs.mousePos.X) {
			rs.totalMouseCapture++
			tail := rs.head.prev
			tail.AddToTail()
			rs.GenerateMouse()
		}
	} else {
		rs.head.move()
	}
	rs.check_game = !rs.check_game
}

func (rs *RetroSnake) Update(delta int64) {

	switch rs.state {
	case STATE_INIT_NEW_GAME:
		rs.initGame()
	case STATE_GAME_RUNNING:
		rs.CalculatePlayTime(delta)

		mx, my := input.Current().GetPosition()
		if mx != -1 && my != -1 {
			if rs.OnClick(mx, my) {
				if rs.DelayElapsed(delta) {
					rs.SetDelay(ANIM_DELAY)
				}
				rs.gamemove()
			}
		} else if rs.DelayElapsed(delta) {
			rs.SetDelay(ANIM_DELAY)
			rs.gamemove()
		}

	case STATE_GAME_PAUSED:
	case STATE_GAME_OVER:
	case STATE_GAME_OVER_HALT:
	}
}

func (rs *RetroSnake) DrawImageAt(screen *ebiten.Image, img *ebiten.Image, x, y int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, op)
}

func (rs *RetroSnake) Draw(screen *ebiten.Image) {
	rs.DrawBoard(screen)
	rs.DrawMouse(screen)
	rs.DrawSnake(screen)
	rs.DrawControlPanel(screen)
	rs.DrawTimer(screen)

}

func (rs *RetroSnake) DrawBoard(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(0), float64(0))
	screen.DrawImage(rs.bgImg, op)
}

func (rs *RetroSnake) DrawMouse(screen *ebiten.Image) {
	if rs.mousePos.X != -1 {
		x := rs.playAreaX + rs.mousePos.X*CELL_WIDTH
		y := rs.playAreaY + rs.mousePos.Y*CELL_HEIGHT
		rs.DrawImageAt(screen, rs.mouseImg, x, y)
	}
}

func (rs *RetroSnake) DrawSnake(screen *ebiten.Image) {

	parse := rs.head

	for parse != nil {

		if parse.curCol >= 0 && parse.curCol <= rs.maxCol && parse.curRow >= 0 && parse.curRow <= rs.maxRow {
			x := rs.playAreaX + parse.curCol*CELL_WIDTH
			y := rs.playAreaY + parse.curRow*CELL_HEIGHT

			if parse.partType == SNAKE_PART_HEAD {
				var imghead *ebiten.Image
				switch parse.direction {
				case DIRECTION_UP:
					imghead = rs.headUp
				case DIRECTION_DOWN:
					imghead = rs.headDown
				case DIRECTION_LEFT:
					imghead = rs.headLeft

				case DIRECTION_RIGHT:
					imghead = rs.headRight

				}

				rs.DrawImageAt(screen, imghead, x, y)

			} else {
				if !(parse.curRow == rs.head.curRow && parse.curCol == rs.head.curCol) {
					if parse.partType == SNAKE_PART_BODY {
						rs.DrawImageAt(screen, rs.greenRound, x, y)
					} else if parse.partType == SNAKE_PART_TAIL {
						rs.DrawImageAt(screen, rs.yellowRound, x, y)
					}
				}
			}

		}

		if parse.partType == SNAKE_PART_TAIL {
			break
		}
		parse = parse.next
	}
}

func (rs *RetroSnake) DrawControlPanel(screen *ebiten.Image) {

	rs.DrawImageAt(screen, rs.upArrow, rs.upButtonX, rs.upButtonY)

	rs.DrawImageAt(screen, rs.downArrow, rs.downButtonX, rs.downButtonY)

	rs.DrawImageAt(screen, rs.leftArrow, rs.leftButtonX, rs.leftButtonY)

	rs.DrawImageAt(screen, rs.rightArrow, rs.rightButtonX, rs.rightButtonY)

	if rs.isReverse {
		rs.DrawImageAt(screen, rs.reverseSel, rs.reverseButtonX, rs.reverseButtonY)

	} else {
		rs.DrawImageAt(screen, rs.reverseDesl, rs.reverseButtonX, rs.reverseButtonY)
	}

	if rs.isPaused {
		rs.DrawImageAt(screen, rs.resumeButton, rs.pauseButtonX, rs.pauseButtonY)

	} else {
		rs.DrawImageAt(screen, rs.pauseButton, rs.pauseButtonX, rs.pauseButtonY)
	}
	rs.DrawImageAt(screen, rs.newButton, rs.newButtonX, rs.newButtonY)
}

func (rs *RetroSnake) DrawTimer(screen *ebiten.Image) {

	rs.DrawImageAt(screen, rs.timeBg, rs.timerBgX, rs.timerBgY)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(rs.timerStrX), float64(rs.timerStrY))
	text.DrawWithOptions(screen, rs.txtTimer, rs.gameFont, op)
}

func (r *RetroSnake) LoadImages() {
	var err error = nil
	r.about, err = r.rm.LoadImage("about.png")
	if err != nil {
		log.Fatal(err)
	}
	r.bgImg, err = r.rm.LoadImage("bg.png")
	if err != nil {
		log.Fatal(err)
	}

	r.downArrow, err = r.rm.LoadImage("down_arrow.png")
	if err != nil {
		log.Fatal(err)
	}

	r.greenRound, err = r.rm.LoadImage("green_round.png")
	if err != nil {
		log.Fatal(err)
	}

	r.headDown, err = r.rm.LoadImage("head_down.png")
	if err != nil {
		log.Fatal(err)
	}

	r.headLeft, err = r.rm.LoadImage("head_left.png")
	if err != nil {
		log.Fatal(err)
	}

	r.headRight, err = r.rm.LoadImage("head_right.png")
	if err != nil {
		log.Fatal(err)
	}

	r.headUp, err = r.rm.LoadImage("head_up.png")
	if err != nil {
		log.Fatal(err)
	}

	r.leftArrow, err = r.rm.LoadImage("left_arrow.png")
	if err != nil {
		log.Fatal(err)
	}

	r.mouseImg, err = r.rm.LoadImage("mouse_img.png")
	if err != nil {
		log.Fatal(err)
	}

	r.newButton, err = r.rm.LoadImage("new_button.png")
	if err != nil {
		log.Fatal(err)
	}

	r.pauseButton, err = r.rm.LoadImage("pause_button.png")
	if err != nil {
		log.Fatal(err)
	}

	r.resumeButton, err = r.rm.LoadImage("resume_button.png")
	if err != nil {
		log.Fatal(err)
	}

	r.reverseDesl, err = r.rm.LoadImage("reverse_desl.png")
	if err != nil {
		log.Fatal(err)
	}

	r.reverseSel, err = r.rm.LoadImage("reverse_sel.png")
	if err != nil {
		log.Fatal(err)
	}

	r.rightArrow, err = r.rm.LoadImage("right_arrow.png")
	if err != nil {
		log.Fatal(err)
	}

	r.timeBg, err = r.rm.LoadImage("time_bg.png")
	if err != nil {
		log.Fatal(err)
	}

	r.upArrow, err = r.rm.LoadImage("up_arrow.png")
	if err != nil {
		log.Fatal(err)
	}

	r.yellowRound, err = r.rm.LoadImage("yellow_round.png")
	if err != nil {
		log.Fatal(err)
	}

	r.overlay, err = r.rm.LoadImage("overlay.png")
	if err != nil {
		log.Fatal(err)
	}
}

func (rs *RetroSnake) LoadFonts() {
	var err error

	rs.TOONEYNO_ttf, err = rs.rm.LoadFont("JustAnotherHand.ttf")
	if err != nil {
		log.Fatal(err)
	}

	rs.gameFont, err = opentype.NewFace(rs.TOONEYNO_ttf, &opentype.FaceOptions{
		Size:    30,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal(err)
	}

}

func (rs *RetroSnake) CreateSnake() {

	rs.head = NewSnakePart(SNAKE_PART_HEAD, 0, 0, color.RGBA{0, 0, 0, 0xff}, DIRECTION_RIGHT)

	bp1 := NewSnakePart(SNAKE_PART_BODY, -1, -1, color.RGBA{0, 0, 0, 0xff}, DIRECTION_RIGHT)
	bp1.prev = rs.head
	rs.head.next = bp1

	bp2 := NewSnakePart(SNAKE_PART_BODY, -1, -1, color.RGBA{0, 0, 0, 0xff}, DIRECTION_RIGHT)
	bp2.prev = bp1
	bp1.next = bp2

	tail := NewSnakePart(SNAKE_PART_TAIL, -1, -1, color.RGBA{0, 0, 0, 0xff}, DIRECTION_RIGHT)
	tail.prev = bp2
	bp2.next = tail

	rs.head.prev = tail
	rs.head.SetMaxRowCol(rs.maxRow, rs.maxCol)

}

// TODO: need to modify this logic
func (rs *RetroSnake) GenerateMouse() {

	var row int = 0
	var col int = 0
	var count int = 0

	for {
		row = rs.random.Intn(rs.maxRow)
		col = rs.random.Intn(rs.maxCol)
		var found bool = true
		var parse = rs.head
		for parse != nil {

			if parse.curCol >= 0 && parse.curCol <= rs.maxCol && parse.curRow >= 0 && parse.curRow <= rs.maxRow {
				if parse.curRow == row && parse.curCol == col {
					found = false
					break
				}
			}
			parse = parse.next
		}

		if found {
			break
		} else {
			count++
			if count == 30 {
				break
			}
		}
	}

	log.Printf(" MRow = %d MCol = %d ", row, col)

	rs.mousePos = image.Point{col, row}
}

func NewRetroSnake(rm *ResourceManager, screenwidth int, screenheight int, callback func(int, int64)) *RetroSnake {
	var retrosnake = new(RetroSnake)
	retrosnake.Init(rm, screenwidth, screenheight, callback)
	return retrosnake
}
