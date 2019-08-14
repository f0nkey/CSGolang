package feature

import (
	"F0nkHack/ferror"
	"F0nkHack/memory"
	"F0nkHack/offset"
	"F0nkHack/types"
	"golang.org/x/xerrors"
	"unsafe"
)

var windowSize = types.GetWindowSize()

func GetViewMatrix(editor *memory.Editor) (memory.CSMatrix, error) {
	//viewMatrix4x4 = readProccessMemoryViewMatrix3(dllClient + oViewMatrix)
	sz := int32(unsafe.Sizeof(memory.CSMatrix{}))
	r, err := editor.Read(sz,editor.DLLClient + offset.Signatures.DwViewMatrix)
	if ferror.ErrIsImportant(err) {
		return memory.CSMatrix{}, xerrors.Errorf("getViewMatrix: %w", err)
	}
	return r.CSMatrix(), nil
}


func WorldToScreen(world memory.Vector3, w2s [4][4]float32) memory.Vector2 {
	var res memory.Vector2
	var w float32

	res.X = w2s[0][0]*world.X + w2s[0][1]*world.Y + w2s[0][2]*world.Z + w2s[0][3]
	res.Y = w2s[1][0]*world.X + w2s[1][1]*world.Y + w2s[1][2]*world.Z + w2s[1][3]
	w = w2s[3][0]*world.X + w2s[3][1]*world.Y + w2s[3][2]*world.Z + w2s[3][3]

	if w < 0.001 {
		res.X *= 100000.00
		res.Y *= 100000.00
	} else {
		var invw float32 = 1.000 / w

		res.X *= invw
		res.Y *= invw
	}

	var sw float32 = float32(windowSize.Width) //todo: change to curr window resolution
	var sh float32 = float32(windowSize.Height)

	var x float32 = float32(sw) / 2.0
	var y float32 = float32(sh) / 2.0

	x += 0.5*res.X*float32(sw) + 0.5
	y -= 0.5*res.Y*float32(sh) + 0.5

	res.X = x
	res.Y = y
	return res
}