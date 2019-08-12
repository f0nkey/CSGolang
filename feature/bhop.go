package feature

import (
	"F0nkHack3/ferror"
	"F0nkHack3/memory"
	"F0nkHack3/offset"
	"github.com/AllenDang/w32"
	"log"
	"time"
)

func BHop(mem *memory.Editor){
	dllClient := mem.DLLClient
	//flags := mem.Read(4,dllClient+offset.LocalPlayer, offset.Flags).Uintptr()


	for {
		keySpace := w32.GetAsyncKeyState(0x20)
		if keySpace == 0 {
			continue
		}

		r, err := mem.Read(4,dllClient+offset.Signatures.DwLocalPlayer, offset.Netvars.MFFlags)
		if ferror.ErrIsImportant(err) {
			log.Println("Bhop:",err)
		}
		onGround := r.Uintptr() & (1 << 0)
		if onGround == 1 {
			err := mem.Write(dllClient+offset.Signatures.DwForceJump, int32(6)) // 6 is jump code
			if err != nil {
				log.Println(err)
			}
		}

		time.Sleep(time.Millisecond * 3)
	}

}
