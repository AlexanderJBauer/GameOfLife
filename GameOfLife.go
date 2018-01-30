// Alexander Bauer
// Notes: Lots of code came from here https://golang.org/doc/play/life.go
// and https://github.com/faiface/pixel/, see pixelLICENSE

package main

import (
	"bufio"
	"fmt"
	"image"
	"math"
	"math/rand"
	"os"
	"time"

	"image/color"
	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"./grid"
)

var (
	NUMROWS    = 0
	SQUARESIZE = 5
	OPTION     = "r"
	play       = true
)

// Field represents a two-dimensional field of cells.
type Field struct {
	s    [][]bool
	w, h int
}

// NewField returns an empty field of the specified width and height.
func NewField(w, h int) *Field {
	s := make([][]bool, h)
	for i := range s {
		s[i] = make([]bool, w)
	}
	return &Field{s: s, w: w, h: h}
}

// Set sets the state of the specified cell to the given value.
func (f *Field) Set(x, y int, b bool) {
	f.s[y][x] = b
}

// Alive reports whether the specified cell is alive.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *Field) Alive(x, y int) bool {
	x += f.w
	x %= f.w
	y += f.h
	y %= f.h
	return f.s[y][x]
}

// Next returns the state of the specified cell at the next time step.
func (f *Field) Next(x, y int) bool {
	// Count the adjacent cells that are alive.
	alive := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && f.Alive(x+i, y+j) {
				alive++
			}
		}
	}
	// Return next state according to the game rules:
	//   exactly 3 neighbors: on,
	//   exactly 2 neighbors: maintain current state,
	//   otherwise: off.
	return alive == 3 || alive == 2 && f.Alive(x, y)
}

// Life stores the state of a round of Conway's Game of Life.
type Life struct {
	a, b *Field
	w, h int
}

// NewLife returns a new Life game state with initial state.
func NewLife(w, h int, option string) *Life {
	a := NewField(w, h)
	switch option {
	case "r":
		for i := 0; i < (w * h / 4); i++ {
			a.Set(rand.Intn(w), rand.Intn(h), true)
		}
	case "e":
		a.Set(w/2, h/2+2, true)
		a.Set(w/2, h/2-2, true)
		a.Set(w/2+2, h/2+2, true)
		a.Set(w/2+2, h/2+1, true)
		a.Set(w/2+2, h/2+0, true)
		a.Set(w/2+2, h/2-1, true)
		a.Set(w/2+2, h/2-2, true)
		a.Set(w/2-2, h/2+2, true)
		a.Set(w/2-2, h/2+1, true)
		a.Set(w/2-2, h/2+0, true)
		a.Set(w/2-2, h/2-1, true)
		a.Set(w/2-2, h/2-2, true)
	case "t":
		a.Set(w/2-5, h/2, true)
		a.Set(w/2-4, h/2, true)
		a.Set(w/2-3, h/2, true)
		a.Set(w/2-2, h/2, true)
		a.Set(w/2-1, h/2, true)
		a.Set(w/2+4, h/2, true)
		a.Set(w/2+3, h/2, true)
		a.Set(w/2+2, h/2, true)
		a.Set(w/2+1, h/2, true)
		a.Set(w/2, h/2, true)
	case "g":
		a.Set(w/2+0, h/2+0, true)
		a.Set(w/2-1, h/2+1, true)
		a.Set(w/2-2, h/2+1, true)
		a.Set(w/2-2, h/2+0, true)
		a.Set(w/2-2, h/2-1, true)
		a.Set(w/2-9, h/2+1, true)
		a.Set(w/2-9, h/2+3, true)
		a.Set(w/2-8, h/2+3, true)
		a.Set(w/2-8, h/2+2, true)
		a.Set(w/2-10, h/2+1, true)
		a.Set(w/2-10, h/2+2, true)
		a.Set(w/2-17, h/2+2, true)
		a.Set(w/2-18, h/2+2, true)
		a.Set(w/2-17, h/2+3, true)
		a.Set(w/2-18, h/2+3, true)
		a.Set(w/2+4, h/2+3, true)
		a.Set(w/2+5, h/2+3, true)
		a.Set(w/2+4, h/2+4, true)
		a.Set(w/2+5, h/2+5, true)
		a.Set(w/2+6, h/2+5, true)
		a.Set(w/2+6, h/2+4, true)
		a.Set(w/2+16, h/2+4, true)
		a.Set(w/2+17, h/2+4, true)
		a.Set(w/2+16, h/2+5, true)
		a.Set(w/2+17, h/2+5, true)
		a.Set(w/2+17, h/2-2, true)
		a.Set(w/2+18, h/2-2, true)
		a.Set(w/2+17, h/2-3, true)
		a.Set(w/2+19, h/2-3, true)
		a.Set(w/2+17, h/2-4, true)
		a.Set(w/2+6, h/2-7, true)
		a.Set(w/2+6, h/2-8, true)
		a.Set(w/2+7, h/2-7, true)
		a.Set(w/2+7, h/2-9, true)
		a.Set(w/2+8, h/2-7, true)
	default:
		for i := 0; i < (w * h / 4); i++ {
			a.Set(rand.Intn(w), rand.Intn(h), true)
		}
	}
	return &Life{
		a: a, b: NewField(w, h),
		w: w, h: h,
	}
}

// Step advances the game by one instant, recomputing and updating all cells.
func (l *Life) Step() {
	// Update the state of the next field (b) from the current field (a).
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			l.b.Set(x, y, l.a.Next(x, y))
		}
	}
	// Swap fields a and b.
	l.a, l.b = l.b, l.a
}

// Loads picture from given path name
func loadPicture(path string) (*pixel.PictureData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() //Waits till function is done to make this call
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

// Creates the frame presented in the window
func DrawBackground(tilemap, spritesheet *pixel.PictureData, cells *Life, rows, squaresize int) (*pixelgl.Canvas, int) {
	background := pixelgl.NewCanvas(tilemap.Bounds())
	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	liveCells := 0
	for x := 0; x <= rows; x++ {
		for y := 0; y <= rows; y++ {
			if cells.a.Alive(x, y) {
				block := pixel.NewSprite(spritesheet, spritesheet.Bounds())
				block.Draw(batch, pixel.IM.Moved(pixel.V(float64(x*squaresize-2), float64(y*squaresize-2)))) //Bug here, only works for size 5
				liveCells++
			}
		}
	}
	pixel.NewSprite(tilemap, tilemap.Bounds()).Draw(background, pixel.IM.Moved(background.Bounds().Center()))
	batch.Draw(background)
	return background, liveCells
}

func run() {

	cells := NewLife(NUMROWS, NUMROWS, OPTION)

	grid.MakeGrid(NUMROWS, SQUARESIZE, color.RGBA{0, 0, 0, 255}, "grid")
	grid.MakeGrid(SQUARESIZE-2, 1, color.RGBA{255, 255, 0, 255}, "block")

	config := pixelgl.WindowConfig{
		Title:  "Conway's Game of Life",
		Bounds: pixel.R(0, 0, 800, 800),
		//VSync:  true,
	}
	win, err := pixelgl.NewWindow(config)
	if err != nil {
		panic(err)
	}

	tilemap, err := loadPicture("grid.png")
	if err != nil {
		panic(err)
	}

	spritesheet, err := loadPicture("block.png")
	if err != nil {
		panic(err)
	}

	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
	)

	var (
		totalFrames = 0
	)

	// something strange is happening here
	// play := false

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		if win.Pressed(pixelgl.KeySpace) {
			play = !play
			time.Sleep(100 * time.Millisecond)
		}

		if win.Pressed(pixelgl.KeyN) {
			cells.Step()
			totalFrames++
			time.Sleep(100 * time.Millisecond)
		}

		if win.Pressed(pixelgl.KeyR) {
			cells = NewLife(NUMROWS, NUMROWS, OPTION)
			totalFrames = 0
			time.Sleep(100 * time.Millisecond)
		}

		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
		}
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		if play {
			cells.Step()
			totalFrames++
		}

		win.Clear(colornames.Gray)
		background, liveCells := DrawBackground(tilemap, spritesheet, cells, NUMROWS, SQUARESIZE)
		background.Draw(win, pixel.IM)
		win.Update()

		win.SetTitle(fmt.Sprintf("%s | Steps: %d | Live Tiles: %d", config.Title, totalFrames, liveCells))

	}
}

func main() {

	//reading an integer
	var rows int
	fmt.Println("What is the preferred number of grid rows?")
	fmt.Scan(&rows)

	NUMROWS = rows

	//reading a string
	reader := bufio.NewScanner(os.Stdin)
	var name string
	fmt.Println("Which pattern would you like to start with? The options are: \nRandom\nExploder\nTen Cell Row\nGlider Gun\nPlease enter first letter of preferred option and hit enter")
	reader.Scan()
	name = reader.Text()

	switch name {
	case "R":
		OPTION = "r"
	case "r":
		OPTION = "r"
	case "E":
		OPTION = "e"
	case "e":
		OPTION = "e"
	case "T":
		OPTION = "t"
	case "t":
		OPTION = "t"
	case "G":
		OPTION = "g"
	case "g":
		OPTION = "g"
	default:
		OPTION = "r"
	}

	// Start application
	pixelgl.Run(run)
}
