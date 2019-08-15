package feature

import (
	"F0nkHack/memory"
	"F0nkHack/render"
	"github.com/chewxy/math32"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

func DrawNames(canvas *render.DrawingCanvas, ps *PlayerStore, vm memory.CSMatrix) {

	for _, currPlayer := range ps.Players {

		x,y := textPosition(currPlayer, vm)
		//addLine(0,0,1920,1080,ic.R,120,120)
		canvas.AddText(x,y,currPlayer.Name, colornames.Amber900)

		//addText(fontm, textDisplay, int(basePos.X-(boxWidth/2))-2+1, int(basePos.Y-boxHeight)-3+1, 1, 0, 0)
		//addText(fontm, textDisplay, int(basePos.X-(boxWidth/2))-2, int(basePos.Y-boxHeight)-3, ic.R, ic.G, ic.B)

		//addRect(int(basePos.X-(boxWidth/2)), int(basePos.Y-boxHeight), int(boxWidth), int(boxHeight), ic.R, ic.G, ic.B)
	}
}

func textPosition(currPlayer Player, vm memory.CSMatrix) (int, int) {
	head := currPlayer.BonePositions[Hitbox_Head]
	vecBox := WorldToScreen(memory.Vector3{head.X, head.Y, head.Z + 5}, vm)
	var headPos = memory.Vector2{vecBox.X, vecBox.Y}
	vecBox = WorldToScreen(memory.Vector3{currPlayer.Position.X, currPlayer.Position.Y, currPlayer.Position.Z}, vm)
	var basePos = memory.Vector2{vecBox.X, vecBox.Y}
	var boxHeight float32 = math32.Abs(headPos.Y - basePos.Y)
	var boxWidth = boxHeight / 2

	return int(basePos.X-(boxWidth/2))-2, int(basePos.Y-boxHeight)-3
}


