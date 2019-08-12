package feature

import (
	"F0nkHack3/ferror"
	"F0nkHack3/memory"
	"F0nkHack3/offset"
	"errors"
	"log"
	"unicode"
	"unsafe"
)

type PlayerStore struct {
	Players    []Player
	maxPlayers uintptr
}

func NewPlayerStore(maxPlayers uintptr) *PlayerStore {
	players := make([]Player, maxPlayers)
	return &PlayerStore{players, maxPlayers}
}

func (ps PlayerStore) UpdateAllPlayers(editor *memory.Editor) {
	for i := uintptr(0); i < ps.maxPlayers; i++ {
		currPlayer := Player{}

		playerAddr := editor.DLLClient + offset.Signatures.DwEntityList   + (i * 0x10)
		name := cleanupName(getName(i,playerAddr,editor))
		bones, err := readBones(playerAddr,editor)
		if err != nil {
			log.Println("UpdateAllPlayers:",err) //TODO: remove UpdateAllPlayers from where it doesnt belong
		}

		sz := unsafe.Sizeof(memory.InternalPlayer{})
		r, err := editor.Read(int32(sz), playerAddr, offset.Netvars.MITeamNum)
		if ferror.ErrIsImportant(err){
			log.Println("UpdateAllPlayers:", err)
		}
		pp := r.InternalPlayer() // struct starts at teamNumber


		currPlayer = Player {
			EntListIndex:    i,
			BaseAddr:        editor.DLLClient + offset.Signatures.DwEntityList + (i * 0x10),
			Team:            pp.Team,
			HP:              pp.HP,
			Position:        memory.Vector3{pp.X, pp.Y, pp.Z},
			Name:            name,
			AngleDifference: 0,
			Bones:           bones,
		}

		ps.Players[i] = currPlayer
	}

}

// TODO: verify this actually works
func readBones(playerAddr uintptr, editor *memory.Editor) (map[uintptr]memory.Vector3, error){
	r, err := editor.Read(int32(4), playerAddr + offset.Signatures.MPStudioHdr,0)
	if ferror.ErrIsImportant(err){
		log.Println("readBones:", err)
	}
	studioHdr := r.Int32()
	if studioHdr == 0 {
		return nil, errors.New("readBones: studioHdr was read to be 0")
	}

	r, err = editor.Read(int32(4), playerAddr + offset.Netvars.MDwBoneMatrix)
	if ferror.ErrIsImportant(err){
		log.Println("readBones:", err)
	}
	bonePre := r.Int32()

	r, err = editor.Read(int32(4), uintptr(studioHdr+0xB0))
	if ferror.ErrIsImportant(err){
		log.Println("readBones:", err)
	}
	hitBoxSetIndex := r.Int32()
	if hitBoxSetIndex < 0 {
		return nil, errors.New("Hitbox set index lower than 0")
	}

	var studioHitboxSet = studioHdr + hitBoxSetIndex

	r, err = editor.Read(int32(4), uintptr(studioHitboxSet+0x4))
	if ferror.ErrIsImportant(err){
		log.Println("readBones:", err)
	}
	numHitboxes := r.Int32()
	if numHitboxes <= 0 {
		return nil, errors.New("No num hitboxes present")
	}

	r, err = editor.Read(int32(4), uintptr(studioHitboxSet+0x8))
	if ferror.ErrIsImportant(err){
		log.Println("readBones:", err)
	}
	hitboxIndex := r.Int32()
	var start int32 = 0
	it := start

	Bones := make(map[uintptr]memory.Vector3)
	for it = 0; it < numHitboxes; it++ {

		//studioHdrData := readProccessMemoryBBox(0x44*(it-start) + hitboxIndex + studioHitboxSet)
		r, err = editor.Read(int32(4), uintptr(0x44*(it-start) + hitboxIndex + studioHitboxSet))
		if ferror.ErrIsImportant(err){
			log.Println("readBones:", err)
		}
		studioHdrData := r.HitBox()

		r, err = editor.Read(int32(4), uintptr(bonePre + 0x30*studioHdrData.Bone))
		if ferror.ErrIsImportant(err){
			log.Println("readBones:", err)
		}
		theBone := r.CSMatrix()

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
	return string(newName[30:]) //TODO: Less hacky way to cleanup text, or remove need to cleanup text
}

func getName(playerIndex uintptr, playerAddr uintptr, editor *memory.Editor) string {
	var ECX_DISP uintptr = 0x40
	var EDX_DISP uintptr = 0x0C
	var INFO_OFFSET uintptr = 0x28
	var ENTRY_SIZE uintptr = 0x34

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
	EntListIndex    uintptr
	BaseAddr        uintptr
	Team            int32
	HP              int32
	Position        memory.Vector3
	Name            string
	AngleDifference float32             //from head to center of screen
	Bones           map[uintptr]memory.Vector3 //using Hitbox consts to refer to bones
}

