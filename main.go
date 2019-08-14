package main

import (
	"F0nkHack/feature"
	"F0nkHack/ferror"
	"F0nkHack/memory"
	"F0nkHack/offset"
	"F0nkHack/render"
	"F0nkHack/types"
	"fmt"
	"log"
	"time"
)

var enableVisuals = true



var windowSize = types.GetWindowSize()

func main() {

	offset.MarshalOffsets()
	memEditor := memory.NewEditor("csgo.exe")
	ps := feature.NewPlayerStore(10)
	ps.UpdateAllPlayers(memEditor)

	for _, value := range ps.Players {
		fmt.Println("HP",value.HP,value.Team, value.Position, value.Name)
	}

	fmt.Printf("%#x\n",offset.Signatures.DwLocalPlayer)

	r, err := memEditor.Read(4, memEditor.DLLClient + offset.Signatures.DwLocalPlayer, 0x100)
	if ferror.ErrIsImportant(err) {
		log.Println("main:",err)
	}
	var hp = r.Int32()
	fmt.Println("small", hp)

	go feature.BHop(memEditor)

	if enableVisuals {
		gw := render.NewDrawingOverlay("Counter-Strike: Global Offensive", "F0nkOverlay",
			func(canvas *render.DrawingCanvas) {

				ps.UpdateAllPlayers(memEditor)
				vm, err := feature.GetViewMatrix(memEditor)
				if err != nil {
					log.Println(err)
				}

				feature.DrawBones(canvas,ps,vm, windowSize)
		})
		gw.RunGL()
	}

	for {
		time.Sleep(time.Millisecond * 5)
	}

}


