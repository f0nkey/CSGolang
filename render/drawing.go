package render

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	tx "github.com/faiface/pixel/text"
	"image/color"
)

type DrawingCanvas struct {
	*imdraw.IMDraw
	windowWidth, windowHeight int
	*tx.Atlas
	*pixelgl.Window
}

func (d DrawingCanvas) AddLine(x, y, cx, cy int, thickness float32, color color.RGBA) {
	x, y, cx, cy = fixCoordinates(x, y, cx, cy, d.windowHeight)

	d.IMDraw.Color = color
	d.IMDraw.EndShape = imdraw.SharpEndShape
	d.IMDraw.Push(pixel.V(float64(x), float64(y)), pixel.V(float64(cx), float64(cy)))
	d.IMDraw.Line(float64(thickness))

}

func (d DrawingCanvas) AddText(x, y int, text string, color color.RGBA) {

	x, y, _, _ = fixCoordinates(x,y,0,0, d.windowHeight)
	basicTxt := tx.New(pixel.V(float64(x), float64(y)), d.Atlas)
	basicTxt.Color = color
	fmt.Fprintln(basicTxt, text)
	basicTxt.Draw(d.Window, pixel.IM)

}

// since IMDraw doesn't use coords starting from top left
func fixCoordinates(x, y, cx, cy, windowHeight int) (int, int, int, int) {
	return x, windowHeight - y, cx, windowHeight - cy
}
