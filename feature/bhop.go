package feature

import (
	"CSGolang/ferror"
	"CSGolang/memory"
	"CSGolang/offset"
	"github.com/AllenDang/w32"
	"log"
	"time"
)

func BHop(mem *memory.Editor, activate bool) {
	dllClient := mem.DLLClient

	for {
		if !activate {
			time.Sleep(time.Millisecond * 3)
			continue
		}
		keySpace := w32.GetAsyncKeyState(0x20)
		if keySpace == 0 {
			continue
		}

		r, err := mem.Read(4, dllClient+offset.Signatures.DwLocalPlayer, offset.Netvars.MFFlags)
		if ferror.ErrIsImportant(err) {
			log.Println("BHop:", err)
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
