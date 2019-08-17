package main

import (
	"F0nkHack/feature"
	"F0nkHack/memory"
	"F0nkHack/offset"
	"F0nkHack/render"
	"F0nkHack/types"
	"log"
	"time"
)

var enableVisuals = true

var windowSize = types.GetWindowSize()

func main() {

	config := feature.InitConfig()
	go feature.WatchConfig(config)
	go feature.RunConfigEndpoint(config)

	offset.InitOffsets()
	memEditor := memory.NewEditor("csgo.exe")

	ps := feature.NewPlayerStore(10)
	ps.UpdateAllPlayers(memEditor)

	go feature.BHop(memEditor, config.Toggles.Bhop)

	gw := render.NewDrawingOverlay("Counter-Strike: Global Offensive", "F0nkOverlay",
		func(canvas *render.DrawingCanvas) {

			ps.UpdateAllPlayers(memEditor)
			vm, err := feature.GetViewMatrix(memEditor)
			if err != nil {
				log.Println(err)
			}

			if config.Toggles.Skeleton {
				feature.DrawBones(canvas, ps, vm, windowSize, config.ColorModes.Skeleton)
			}

			if config.Toggles.Name {
				feature.DrawNames(canvas, ps, vm, config.ColorModes.Name)
			}

		})
	gw.RunGL()

	for {
		time.Sleep(time.Millisecond * 5)
	}

}
