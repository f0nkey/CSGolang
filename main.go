package main

import (
	"F0nkHack3/feature"
	"F0nkHack3/ferror"
	"F0nkHack3/memory"
	"F0nkHack3/offset"
	"F0nkHack3/render"
	"fmt"
	"log"
	"time"
)

var enableVisuals = false

func main() {
	offset.MarshalOffsets()
	memEditor := memory.NewEditor("csgo.exe")
	ps := feature.NewPlayerStore(10)
	ps.UpdateAllPlayers(memEditor)

	//fmt.Println("TEAMS", ps.Players[0].HP)
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

	//var mm = 1
	if enableVisuals {
		gw := render.NewDrawingOverlay("Counter-Strike: Global Offensive", "F0nkOverlay",
			func(canvas *render.DrawingCanvas) {

			//canvas.AddLine(50, mm, 900, 900, 1, colornames.Amber900)
			//mm++

		})
		gw.RunGL()
	}

	for {
		time.Sleep(time.Millisecond * 5)
	}


}


