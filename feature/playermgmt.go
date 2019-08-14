package feature

import (
	"F0nkHack/ferror"
	"F0nkHack/memory"
	"F0nkHack/offset"
	"fmt"
	"golang.org/x/xerrors"
	"log"
	"unicode"
	"unsafe"
)

type PlayerStore struct {
	Players    []Player
	maxPlayers int32
}

func NewPlayerStore(maxPlayers int32) *PlayerStore {
	players := make([]Player, maxPlayers)
	return &PlayerStore{players, maxPlayers}
}

func (ps PlayerStore) UpdateAllPlayers(editor *memory.Editor) {
	var dcount int
	for i := int32(0); i < ps.maxPlayers; i++ {
		currPlayer := Player{}

		playerAddr := editor.DLLClient + offset.Signatures.DwEntityList + (i * 0x10)


		sz := unsafe.Sizeof(memory.InternalPlayer{})
		r, err := editor.Read(int32(sz), playerAddr, offset.Netvars.MITeamNum) //struct starts at teamNumber
		if ferror.ErrIsImportant(err){
			log.Println("UpdateAllPlayers:", err)
		}
		pp := r.InternalPlayer()
		if (pp.Team != 3 && pp.Team != 2) || pp.HP <= 0 { // is not a player, or alive
			dcount++
			continue
		}

		bonePositions,err := readBonesWorldPos(playerAddr, editor) //todo: fix me fucker
		if err != nil {
			log.Println("UpdateAllPlayers:",err) //TODO: remove UpdateAllPlayers from where it doesnt belong
		}
		name := cleanupName(getName(i,playerAddr,editor))

		currPlayer = Player {
			EntListIndex:    i,
			BaseAddr:        editor.DLLClient + offset.Signatures.DwEntityList + (i * 0x10),
			Team:            pp.Team,
			HP:              pp.HP,
			Position:        memory.Vector3{pp.X, pp.Y, pp.Z},
			Name:            name,
			AngleDifference: 0,
			BonePositions:   bonePositions,
		}

		ps.Players[i] = currPlayer
	}
	fmt.Println("dcount",dcount)
}


func readBonesWorldPos(playerAddr int32, editor *memory.Editor) (map[uintptr]memory.Vector3, error){
	playerAddr = editor.Read2(4,playerAddr).Int32()
	studioHdr := editor.Read2(4,playerAddr+offset.Signatures.MPStudioHdr).Int32()
	bonePre := editor.Read2(4,playerAddr+offset.Netvars.MDwBoneMatrix).Int32()
	studioHdr = editor.Read2(4,studioHdr).Int32()

	if studioHdr == 0 {
		return nil, xerrors.New("readBonesWorldPos: no character present")
	}

	hitBoxSetIndex := editor.Read2(4,studioHdr+0xB0).Int32()
	if hitBoxSetIndex < 0 {
		return nil, xerrors.New("readBonesWorldPos: Hitbox set index lower than 0")
	}

	var studioHitboxSet = int32(studioHdr + hitBoxSetIndex)
	numHitboxes := editor.Read2(4,studioHitboxSet+0x4).Int32()
	if numHitboxes <= 0 {
		return nil, xerrors.New("readBonesWorldPos: no hitboxes present")
	}

	var hitboxIndex = editor.Read2(4,studioHitboxSet+0x8).Int32()
	var start int32 = 0
	it := start

	var Bones = make(map[uintptr]memory.Vector3)
	for it = 0; it < numHitboxes; it++ {
		sz := int32(unsafe.Sizeof(memory.HitBox{}))
		studioHdrData := editor.Read2(sz,0x44*(it-start) + hitboxIndex + studioHitboxSet).HitBox()

		sz = int32(unsafe.Sizeof(memory.CSMatrix{}))
		theBone := editor.Read2(sz,bonePre + 0x30*studioHdrData.Bone).CSMatrix()

		var bonePos = memory.Vector3{theBone[0][3], theBone[1][3], theBone[2][3]}

		Bones[uintptr(it)] = bonePos
	}
	return Bones, nil
}

func cleanupName(s string) string {
	runes := []rune(s)
	newName := make([]rune,len(runes))
	for _, char := range s {
		if unicode.IsPrint(char) {
			newName = append(newName,char)
		}
	}
	return string(newName[30:]) //TODO: Less hacky way to cleanup text, or remove need to cleanup text with better getName
}

func getName(playerIndex, playerAddr int32, editor *memory.Editor) string {
	var ECX_DISP int32 = 0x40
	var EDX_DISP int32 = 0x0C
	var INFO_OFFSET int32 = 0x28
	var ENTRY_SIZE int32 = 0x34

	r, err := editor.Read(int32(32), editor.DLLEngine+offset.Signatures.DwClientState,
		offset.Signatures.DwClientStatePlayerInfo,
		ECX_DISP,
		EDX_DISP,
		INFO_OFFSET+playerIndex*ENTRY_SIZE,0)

	if ferror.ErrIsImportant(err){
		log.Println("getName:",err)
	}

	return r.String()
}

type Player struct {
	EntListIndex    int32
	BaseAddr        int32
	Team            int32
	HP              int32
	Position        memory.Vector3
	Name            string
	AngleDifference float32                    //from head to center of screen
	BonePositions   map[uintptr]memory.Vector3 //using Hitbox consts to refer to bones
}

//exp funcs

//var (
//	modkernel32            = syscall.NewLazyDLL("kernel32.dll")
//	procWriteProcessMemory = modkernel32.NewProc("WriteProcessMemory")
//	procReadProcessMemory  = modkernel32.NewProc("ReadProcessMemory")
//)
//
//var csgoPID, err = memory.FindProcessByName("csgo.exe") //using to get the PID of csgo.exe
//var handleCSGO, _ = w32.OpenProcess(w32.PROCESS_ALL_ACCESS, false, uintptr(csgoPID))
//var dllClient = uintptr(unsafe.Pointer(memory.GetDLLModuleAddress("client_panorama.dll", csgoPID)))
//var dllEngine = uintptr(unsafe.Pointer(memory.GetDLLModuleAddress("engine.dll", csgoPID)))
//
//func readBonesExp(BaseAddr uintptr) map[uintptr]memory.Vector3 {
//	BaseAddr = readProccessMemory(BaseAddr,4)
//
//	studioHdr := readProccessMemory(BaseAddr+uintptr(offset.Signatures.MPStudioHdr), 4)
//
//	bonePre := readProccessMemory(BaseAddr+uintptr(offset.Netvars.MDwBoneMatrix), 4)
//
//	hp := readProccessMemory(BaseAddr+0x100, 4)
//	fmt.Println("hp",hp,BaseAddr + 0x100)
//	studioHdr = readProccessMemory(studioHdr, 4)
//	fmt.Println("FMexp=========",studioHdr)
//	if studioHdr == 0 {
//		fmt.Println("Did not populate bonemap:" )
//		fmt.Println("Nothing present")
//		return nil
//	}
//
//	hitBoxSetIndex := readProccessMemory(studioHdr+0xB0, 4)
//	if hitBoxSetIndex < 0 {
//		fmt.Println("Did not populate bonemap:" )
//		fmt.Println("Hitbox set index lower than 0")
//		return nil
//	}
//
//	var studioHitboxSet uintptr = studioHdr + hitBoxSetIndex
//	numHitboxes := readProccessMemory(studioHitboxSet+0x4, 4)
//	if numHitboxes <= 0 {
//		fmt.Println("Did not populate bonemap:")
//		fmt.Println("No num hitboxes")
//		return nil
//	}
//
//	var hitboxIndex uintptr = readProccessMemory(studioHitboxSet+0x8, 4)
//	var start uintptr = 0
//	it := start
//
//
//	//for b := 0; b< 30; b++ {
//	//
//	//	studioHdrData := readProccessMemoryBBox(0x44*(it-start) + hitboxIndex + studioHitboxSet)
//	//
//	//}
//
//
//	var BonePositions = make(map[uintptr]memory.Vector3)
//	for it = 0; it < numHitboxes; it++ {
//
//
//
//
//		studioHdrData := readProccessMemoryBBox(0x44*(it-start) + hitboxIndex + studioHitboxSet)
//		//debugPrintHigh(now8,*p,it,"wWTF PLEASE BWORK",studioHdrData)
//		//debugPrintHigh(now,p,it,"studioHDR",studioHdrData)
//
//		//now2 := time.Now()
//		theBone := readProccessMemoryViewMatrix3(bonePre + 0x30*uintptr(studioHdrData.Bone))
//		//debugPrintHigh(now2,p,it,"bone",theBone)
//
//
//
//		var bonePos = memory.Vector3{theBone[0][3], theBone[1][3], theBone[2][3]}
//		//fmt.Println("StudioBone:",studioHdrData.Bone,bonePos)
//
//		BonePositions[uintptr(it)] = bonePos //using hitBoxIDs array to make a named key instead of creating the names at boneMap init
//	}
//	return BonePositions
//}
//
//func readProccessMemory(lpBaseAddress uintptr, nSize uintptr) uintptr {
//
//	var nBytesRead int
//	buf := make([]byte, nSize)
//	procReadProcessMemory.Call(
//		uintptr(handleCSGO),
//		lpBaseAddress,
//		uintptr(unsafe.Pointer(&buf[0])),
//		nSize,
//		uintptr(unsafe.Pointer(&nBytesRead)),
//	)
//	return uintptr(binary.LittleEndian.Uint32(buf))
//}
//
//func readProccessMemoryBBox(lpBaseAddress uintptr) (viewM BBox) {
//	var nBytesRead int
//	var nSize = unsafe.Sizeof(viewM)
//	buf := make([]byte, nSize)
//
//	_,_,_ = procReadProcessMemory.Call(
//		uintptr(handleCSGO),
//		lpBaseAddress,
//		uintptr(unsafe.Pointer(&buf[0])),
//		nSize,
//		uintptr(unsafe.Pointer(&nBytesRead)),
//	)
//
//
//
//	newReader := bytes.NewBuffer(buf)
//	binary.Read(newReader, binary.LittleEndian, &viewM)
//	//fmt.Printf("%.5f\n",viewM)
//	return viewM
//}
//
//func readProccessMemoryViewMatrix3(lpBaseAddress uintptr) (viewM [4][4]float32) {
//	var nBytesRead int
//	var nSize = unsafe.Sizeof(viewM)
//	buf := make([]byte, nSize)
//	procReadProcessMemory.Call(
//		uintptr(handleCSGO),
//		lpBaseAddress,
//		uintptr(unsafe.Pointer(&buf[0])),
//		nSize,
//		uintptr(unsafe.Pointer(&nBytesRead)),
//	)
//	newReader := bytes.NewBuffer(buf)
//	binary.Read(newReader, binary.LittleEndian, &viewM)
//	//fmt.Printf("%.5f\n",viewM)
//	return viewM
//}
//
//type BBox struct {
//	Bone     int32      /// 0x00
//	Group    int32      /// 0x04
//	Mins     memory.Vector3    /// 0x08
//	Maxs     memory.Vector3    /// 0x14
//	HitBoxID int32      /// 0x20
//	Pad1     [0xC]byte  /// 0x24
//	Radius   float32    /// 0x30
//	Pad2     [0x10]byte /// 0x34
//}

