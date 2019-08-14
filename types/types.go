package types

import (
	"F0nkHack/memory"
	"fmt"
	"github.com/AllenDang/w32"
)

type WindowRect struct {
	Width, Height int32
}

func GetWindowSize() WindowRect {
	w, err := memory.FindWindow("Counter-Strike: Global Offensive")
	if err != nil {
		panic("Did not find Counter-Strike: Global Offensive window. Is CS:GO is open?")
	}
	rect := w32.GetWindowRect(w32.HWND(w))
	fmt.Println("rect",rect)
	return WindowRect{rect.Right - rect.Left,rect.Bottom - rect.Top}
}
