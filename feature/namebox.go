package feature

import (
	"F0nkHack/memory"
	"F0nkHack/render"
	"github.com/chewxy/math32"
	"image/color"
)

func DrawNames(canvas *render.DrawingCanvas, ps *PlayerStore, vm memory.CSMatrix, colorMode int) {

	for _, currPlayer := range ps.Players {
		if currPlayer.IsDormant {
			continue
		}
		x, y := textPosition(currPlayer, vm)
		col := getColor(colorMode, currPlayer.HP, currPlayer.Team)
		canvas.AddText(x+1, y+1, currPlayer.Name, color.RGBA{0, 0, 0, 255})
		canvas.AddText(x, y, currPlayer.Name, col)
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

	return int(basePos.X-(boxWidth/2)) - 2, int(basePos.Y-boxHeight) - 3
}
