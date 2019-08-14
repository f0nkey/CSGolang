package render

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"image/color"
)

type DrawingCanvas struct {
	*imdraw.IMDraw
	windowWidth, windowHeight int
}



func (d DrawingCanvas) AddLine(x,y,cx,cy int, thickness float32, color color.RGBA){
	x,y,cx,cy = fixCoordinates(x,y,cx,cy, d.windowHeight)

	d.IMDraw.Color = color
	d.IMDraw.EndShape = imdraw.SharpEndShape
	d.IMDraw.Push(pixel.V(float64(x), float64(y)), pixel.V(float64(cx), float64(cy)))
	d.IMDraw.Line(float64(thickness))

}

// IMDraw doesn't use opengl coords starting from top left
func fixCoordinates(x,y,cx,cy, windowHeight int)(int, int, int,int) {
	return x,windowHeight-y,cx,windowHeight-cy
}