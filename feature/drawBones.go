package feature

import (
	"CSGolang/memory"
	"CSGolang/render"
	"CSGolang/types"
	"image/color"
)

const (
	Hitbox_Head            = iota
	Hitbox_Neck            = iota
	Hitbox_Pelvis          = iota
	Hitbox_Stomach         = iota
	Hitbox_Chest           = iota
	Hitbox_Lower_Chest     = iota
	Hitbox_Upper_Chest     = iota
	Hitbox_Right_Thigh     = iota
	Hitbox_Left_Thigh      = iota
	Hitbox_Right_Calf      = iota
	Hitbox_Left_Calf       = iota
	Hitbox_Right_Foot      = iota
	Hitbox_Left_Foot       = iota
	Hitbox_Right_Hand      = iota
	Hitbox_Left_Hand       = iota
	Hitbox_Right_Upper_Arm = iota
	Hitbox_Right_Lower_Arm = iota
	Hitbox_Left_Upper_Arm  = iota
	Hitbox_Left_Lower_Arm  = iota
	Hitbox_Last            = iota
)

var hitBoxIDS []string = []string{
	"Hitbox_Head",
	"Hitbox_Neck",
	"Hitbox_Pelvis",
	"Hitbox_Stomach",
	"Hitbox_Chest",
	"Hitbox_Lower_Chest",
	"Hitbox_Upper_Chest",
	"Hitbox_Right_Thigh",
	"Hitbox_Left_Thigh",
	"Hitbox_Right_Calf",
	"Hitbox_Left_Calf",
	"Hitbox_Right_Foot",
	"Hitbox_Left_Foot",
	"Hitbox_Right_Hand",
	"Hitbox_Left_Hand",
	"Hitbox_Right_Upper_Arm",
	"Hitbox_Right_Lower_Arm",
	"Hitbox_Left_Upper_Arm",
	"Hitbox_Left_Lower_Arm",
	"Hitbox_Last",
}

func DrawBones(canvas *render.DrawingCanvas, ps *PlayerStore, viewMatrix memory.CSMatrix, windowSize types.WindowRect, colorMode int) {
	players := ps.Players
	for _, currPlayer := range players {
		if !isWithinFov(currPlayer.Position, viewMatrix, windowSize) || currPlayer.EntListIndex == 0 || currPlayer.IsDormant {
			continue
		}

		var boneLinks = [][]int{
			[]int{Hitbox_Head, Hitbox_Neck, Hitbox_Upper_Chest, Hitbox_Lower_Chest, Hitbox_Chest, Hitbox_Stomach, Hitbox_Pelvis}, // spine
			[]int{Hitbox_Left_Foot, Hitbox_Left_Calf, Hitbox_Left_Thigh, Hitbox_Pelvis},                                          // left leg
			[]int{Hitbox_Right_Foot, Hitbox_Right_Calf, Hitbox_Right_Thigh, Hitbox_Pelvis},                                       // right leg
			[]int{Hitbox_Left_Hand, Hitbox_Left_Lower_Arm, Hitbox_Left_Upper_Arm, Hitbox_Upper_Chest},                            // left arm
			[]int{Hitbox_Right_Hand, Hitbox_Right_Lower_Arm, Hitbox_Right_Upper_Arm, Hitbox_Upper_Chest},                         //right arm
		}

		c := getColor(colorMode, currPlayer.HP, currPlayer.Team)
		for _, boneLink := range boneLinks {
			drawBoneLinks(boneLink, currPlayer.BonePositions, currPlayer, canvas, viewMatrix, c)
		}
	}
}

// TODO: account for player's width
func isWithinFov(playerPos memory.Vector3, viewMatrix memory.CSMatrix, windowSize types.WindowRect) bool {
	feetLevel := WorldToScreen(playerPos, viewMatrix)
	ad := memory.Vector3{playerPos.X, playerPos.Y, playerPos.Z}
	ad.Z += 64 //height of a csgo player
	headLevel := WorldToScreen(ad, viewMatrix)

	a := (int32(feetLevel.X) > windowSize.Width || feetLevel.X < 0) || (int32(feetLevel.Y) > windowSize.Height || feetLevel.Y < 0)
	b := (int32(headLevel.X) > windowSize.Width || headLevel.X < 0) || (int32(headLevel.Y) > windowSize.Height || headLevel.Y < 0)
	if a || b {
		return false
	}
	return true
}

func drawBoneLinks(bones []int, boneMap map[uintptr]memory.Vector3, currPlayer Player, canvas *render.DrawingCanvas, viewMatrix memory.CSMatrix, color color.RGBA) {
	var a uintptr = 0
	var b uintptr = 1

	for b < uintptr(len(bones)) {
		drawBone(boneMap, uintptr(bones[a]), uintptr(bones[b]), currPlayer, canvas, viewMatrix, color)
		a++
		b++
	}
}

//todo: make a DrawBoneOptions struct
func drawBone(boneMap map[uintptr]memory.Vector3, from uintptr, to uintptr, currPlayer Player, canvas *render.DrawingCanvas, viewMatrix memory.CSMatrix, color color.RGBA) {
	vecScreen := WorldToScreen(boneMap[from], viewMatrix)
	vecScreenc := WorldToScreen(boneMap[to], viewMatrix)
	canvas.AddLine(int(vecScreen.X), int(vecScreen.Y), int(vecScreenc.X), int(vecScreenc.Y), 1, color)
}
