package render

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"image/color"
)

type DrawingCanvas struct {
	*imdraw.IMDraw
}

func (d DrawingCanvas) AddLine(x,y,cx,cy int, thickness float32, color color.RGBA){
	d.IMDraw.Color = color
	d.IMDraw.EndShape = imdraw.SharpEndShape
	d.IMDraw.Push(pixel.V(float64(x), float64(y)), pixel.V(float64(cx), float64(cy)))
	d.IMDraw.Line(float64(thickness))

}