package main

import (
	"F0nkHack/feature"
	"F0nkHack/ferror"
	"F0nkHack/memory"
	"F0nkHack/offset"
	"F0nkHack/render"
	"fmt"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"log"
	"net"
	"net/http"
	"time"
)

var enableVisuals = true

var windowWidth, windowHeight float32

func main() {
	initExpVarServer()
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

	var mm = 1
	if enableVisuals {
		gw := render.NewDrawingOverlay("Counter-Strike: Global Offensive", "F0nkOverlay",
			func(canvas *render.DrawingCanvas) {

				ps.UpdateAllPlayers(memEditor)
				vm, err := feature.GetViewMatrix(memEditor)
				if err != nil {
					log.Println(err)
				}

				feature.DrawBones(canvas,ps,vm)
				canvas.AddLine(500, 50, 900, 900, 1, colornames.Amber900)
				mm++

		})
		gw.RunGL()
	}

	for {
		time.Sleep(time.Millisecond * 5)
	}


}

func initExpVarServer() error {
	sock, err := net.Listen("tcp", "localhost:8181")
	if err != nil {
		return err
	}
	go func() {
		fmt.Println("HTTP now available at port 8123")
		http.Serve(sock, nil)
	}()
	return nil
}


