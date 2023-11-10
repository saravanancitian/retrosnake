package retrosnake

import (
	"image/color"
)

const (
	DIRECTION_UP    int = 0
	DIRECTION_DOWN  int = 1
	DIRECTION_LEFT  int = 2
	DIRECTION_RIGHT int = 3
)

const (
	SNAKE_PART_HEAD int = 0
	SNAKE_PART_BODY int = 1
	SNAKE_PART_TAIL int = 2

	PART_WIDTH  int = 30
	PART_HEIGHT int = 30
)

type SnakePart struct {
	partType  int
	color     color.RGBA
	direction int
	prev      *SnakePart
	next      *SnakePart

	maxRow int
	maxCol int

	curRow int
	curCol int

	diff int
	x    int
	y    int
}

func (sp *SnakePart) Init(partType int, row int, col int, color color.RGBA, direction int) {
	sp.partType = partType
	sp.curRow = row
	sp.curCol = col
	sp.color = color
	sp.direction = direction
	sp.diff = 30
	sp.x = 0
	sp.y = 0
}

func (sp *SnakePart) SetMaxRowCol(maxRow, maxCol int) {
	sp.maxRow = maxRow
	sp.maxCol = maxCol
}

func (sp *SnakePart) copyPrev() {
	if sp.prev != nil {
		sp.curRow = sp.prev.curRow
		sp.curCol = sp.prev.curCol
		sp.direction = sp.prev.direction
	}
}

func (sp *SnakePart) AddToTail() {
	if sp.partType == SNAKE_PART_TAIL {
		var bp = NewSnakePart(SNAKE_PART_BODY, sp.curRow, sp.curCol, sp.color, sp.direction)
		bp.prev = sp.prev
		bp.next = sp
		sp.prev.next = bp
		sp.prev = bp
		switch sp.direction {
		case DIRECTION_UP:
			sp.curRow += 1
			// sp.y += PART_HEIGHT;
		case DIRECTION_DOWN:
			sp.curRow -= 1
			//sp.y -= PART_HEIGHT;
		case DIRECTION_LEFT:
			sp.curCol += 1
			//sp.x += PART_WIDTH;
		case DIRECTION_RIGHT:
			sp.curCol -= 1
			//sp.x -= PART_WIDTH;
		}

	}
}

func (sp *SnakePart) CheckCaptureMouse(rowIdx, colIdx int) bool {
	var retVal bool = false
	if sp.partType == SNAKE_PART_HEAD {
		// log.Printf("mr = %v, mc = %v , hr = %v, hc = %v", rowIdx, colIdx, sp.curRow, sp.curCol)
		retVal = (rowIdx == sp.curRow && colIdx == sp.curCol)
	}
	return retVal
}

func (sp *SnakePart) CheckTailBite() bool {
	var retVal bool = false
	if sp.partType == SNAKE_PART_HEAD {
		var parse = sp.next
		for parse != nil {
			retVal = (sp.curRow == parse.curRow && sp.curCol == parse.curCol)
			if retVal {
				break
			}
			parse = parse.next
		}

	}
	return retVal
}

func (sp *SnakePart) changeDirection(direction int) bool {
	var handled bool = false
	if sp.partType == SNAKE_PART_HEAD {
		if ((sp.direction == DIRECTION_LEFT || sp.direction == DIRECTION_RIGHT) && (direction != DIRECTION_LEFT && direction != DIRECTION_RIGHT)) || ((sp.direction == DIRECTION_UP || sp.direction == DIRECTION_DOWN) && (direction != DIRECTION_UP && direction != DIRECTION_DOWN)) {
			var parse = sp.prev
			for parse != sp {
				parse.copyPrev()
				parse = parse.prev
			}
			sp.direction = direction
			sp.moveCoord()
			handled = true
		}
	}

	return handled
}

func (sp *SnakePart) moveCoord() {
	switch sp.direction {
	case DIRECTION_UP:

		sp.curRow -= 1
		if sp.curRow < 0 {
			sp.curRow = sp.maxRow
		}

	case DIRECTION_DOWN:
		sp.curRow += 1
		if sp.curRow > sp.maxRow {
			sp.curRow = 0
		}
	case DIRECTION_LEFT:
		sp.curCol -= 1
		if sp.curCol < 0 {
			sp.curCol = sp.maxCol
		}
	case DIRECTION_RIGHT:
		sp.curCol += 1
		if sp.curCol > sp.maxCol {
			sp.curCol = 0
		}
	}
}

func (sp *SnakePart) move() {
	if sp.partType == SNAKE_PART_HEAD {
		
		var tail = sp.prev
		var parse = tail
		for parse != sp {
			parse.copyPrev()
			parse = parse.prev
		}
		sp.moveCoord()
	}
}

func (sp *SnakePart) copyNext() {
	if sp.next != nil {
		sp.curRow = sp.next.curRow
		sp.curCol = sp.next.curCol
		sp.direction = sp.next.direction
	}
}

func NewSnakePart(partType int, row int, col int, color color.RGBA, direction int) *SnakePart {
	var snakePart = new(SnakePart)
	snakePart.Init(partType, row, col, color, direction)
	return snakePart
}
