package render

import (
	"F0nkHack/memory"
	"fmt"
	"github.com/AllenDang/w32"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"log"
	"time"
)

type guiWindow struct {
	TargetWindow string
	WindowName string
	RenderLoop func(*DrawingCanvas)
	targetWindowRect *w32.RECT
}

func NewDrawingOverlay(targetWindow string, name string, drawLoop func(canvas *DrawingCanvas)) *guiWindow {
	t, err := memory.FindWindow(targetWindow)
	if err != nil {
		log.Println("createWindowAndRenderLoop: error finding target window, ",err)
	}

	return &guiWindow{targetWindow,name,drawLoop,w32.GetWindowRect(w32.HWND(t))}
}

func (g guiWindow) createWindowAndRenderLoop() {

	 fmt.Println(pixel.R(float64(g.targetWindowRect.Left),
		float64(g.targetWindowRect.Top),
		float64(g.targetWindowRect.Right),
		float64(g.targetWindowRect.Bottom)))

	cfg := pixelgl.WindowConfig {
		Title:  g.WindowName,
		Bounds: pixel.R(float64(g.targetWindowRect.Left),
			float64(g.targetWindowRect.Top)-21,
			float64(g.targetWindowRect.Right),
			float64(g.targetWindowRect.Bottom)-21),
		VSync:  false,
		Undecorated: false,
	}
	window, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	canv := DrawingCanvas{imd,int(g.targetWindowRect.Right),int(g.targetWindowRect.Bottom)}

	for !window.Closed() {
		time.Sleep(time.Millisecond * 16)
		window.Clear(color.RGBA{0x00, 0x00, 0x00, 0x00})
		imd.Clear()

		g.RenderLoop(&canv)
		imd.Draw(window)

		window.Update()
	}
}

func (g guiWindow) RunGL() {
	go g.setAlphaWindow()
	pixelgl.Run(g.createWindowAndRenderLoop)
}

func (g guiWindow) setAlphaWindow() {
	time.Sleep(100 * time.Millisecond)
	preHwnd, err := memory.FindWindow(g.WindowName)
	if err != nil {
		fmt.Println(err)
	}
	hwnd := w32.HWND(preHwnd)

	pflags := w32.GetWindowLong(hwnd, w32.GWL_STYLE)
	pflags &= ^(w32.WS_CAPTION | w32.WS_THICKFRAME | w32.WS_MINIMIZE | w32.WS_MAXIMIZE | w32.WS_SYSMENU)
	ok := w32.SetWindowLong(hwnd, w32.GWL_STYLE, uint32(pflags))
	if ok == 0 {
		log.Println("SetWindowLong: Error setting GWL_STYLE ", ok)
	}

	marg := w32.MARGINS{-1, -1, -1, -1}
	hr := w32.DwmExtendFrameIntoClientArea(hwnd, &marg)
	if hr != w32.S_OK {
		log.Println("DwmExtendFrameIntoClientArea: could not set frame ", hr)
	}

	flags := w32.GetWindowLong(hwnd, w32.GWL_EXSTYLE);
	flags = flags&(^(w32.WS_EX_STATICEDGE | w32.WS_EX_CLIENTEDGE));
	ok = w32.SetWindowLong(hwnd, w32.GWL_EXSTYLE, uint32(flags) | w32.WS_EX_TRANSPARENT | w32.WS_EX_NOACTIVATE | w32.WS_EX_LAYERED)
	if ok == 0 {
		log.Println("SetWindowLong: Error setting GWL_EXSTYLE ", ok)
	}

	ok2 := w32.SetWindowPos(hwnd,w32.HWND_TOPMOST,0,0,0,0, w32.SWP_FRAMECHANGED | w32.SWP_NOMOVE | w32.SWP_NOSIZE | w32.SWP_NOREPOSITION | w32.SWP_NOSIZE | w32.SWP_NOACTIVATE)
	if ok2 == false {
		log.Println("SetWindowPos: Error setting flags ", ok)
	}
	log.Println("SetWindowPos:", ok)

	if !ok2 {
		g.setAlphaWindow()
	}
}