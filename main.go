package main

import (
	"CSGolang/feature"
	"CSGolang/memory"
	"CSGolang/offset"
	"CSGolang/render"
	"CSGolang/types"
	"log"
	"time"
)

var enableVisuals = true

var windowSize = types.GetWindowSize()

func main() {
	config := feature.InitConfig()
	go feature.WatchConfig(config)
	go feature.RunConfigEndpoint(config) //localhost:9991
	go feature.ServeWebGUI() //localhost:8085
	//
	offset.InitOffsets()
	memEditor := memory.NewEditor("csgo.exe")

	ps := feature.NewPlayerStore(12)
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
