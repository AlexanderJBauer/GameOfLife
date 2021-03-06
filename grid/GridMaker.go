package grid

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

// hLine draws a horizontal line
func hLine(x1, y, x2 int, col color.Color, img *image.RGBA) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

// vLine draws a veritcal line
func vLine(x, y1, y2 int, col color.Color, img *image.RGBA) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

// rect draws a rectangle utilizing hLine() and vLine()
func rect(x1, y1, x2, y2 int, col color.Color, img *image.RGBA) {
	hLine(x1, y1, x2, col, img)
	hLine(x1, y2, x2, col, img)
	vLine(x1, y1, y2, col, img)
	vLine(x2, y1, y2, col, img)
}

func MakeGrid(numRows, squareSize int, col color.Color, imageName string) {

	var img *image.RGBA

	img = image.NewRGBA(image.Rect(0, 0, numRows*squareSize+1, numRows*squareSize+1))

	for i := 0; i < numRows; i++ {
		for j := 0; j < numRows; j++ {
			rect(squareSize*i, squareSize*j, squareSize*i+squareSize, squareSize*j+squareSize, col, img)
		}
	}

	f, err := os.Create(imageName + ".png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
