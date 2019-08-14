package feature

import (
	"F0nkHack/ferror"
	"F0nkHack/memory"
	"F0nkHack/offset"
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
	for i := int32(0); i < ps.maxPlayers; i++ {
		currPlayer := Player{}

		playerAddr := editor.DLLClient + offset.Signatures.DwEntityList + (i * 0x10)

		sz := unsafe.Sizeof(memory.InternalPlayer{})
		r, err := editor.Read(int32(sz), playerAddr, offset.Netvars.MITeamNum) //struct starts at teamNumber
		if ferror.ErrIsImportant(err) {
			log.Println("UpdateAllPlayers:", err)
		}
		pp := r.InternalPlayer()
		if (pp.Team != 3 && pp.Team != 2) || pp.HP <= 0 { // is not a player, or alive
			continue
		}

		bonePositions, err := readBonesWorldPos(playerAddr, editor) //todo: fix me fucker
		if err != nil {
			log.Println("UpdateAllPlayers:", err) //TODO: remove UpdateAllPlayers from where it doesnt belong
		}
		name := cleanupName(getName(i, playerAddr, editor))

		currPlayer = Player{
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
}

func readBonesWorldPos(playerAddr int32, editor *memory.Editor) (map[uintptr]memory.Vector3, error) {
	playerAddr = editor.Read2(4, playerAddr).Int32()
	studioHdr := editor.Read2(4, playerAddr+offset.Signatures.MPStudioHdr).Int32()
	bonePre := editor.Read2(4, playerAddr+offset.Netvars.MDwBoneMatrix).Int32()
	studioHdr = editor.Read2(4, studioHdr).Int32()

	if studioHdr == 0 {
		return nil, xerrors.New("readBonesWorldPos: no character present")
	}

	hitBoxSetIndex := editor.Read2(4, studioHdr+0xB0).Int32()
	if hitBoxSetIndex < 0 {
		return nil, xerrors.New("readBonesWorldPos: Hitbox set index lower than 0")
	}

	var studioHitboxSet = int32(studioHdr + hitBoxSetIndex)
	numHitboxes := editor.Read2(4, studioHitboxSet+0x4).Int32()
	if numHitboxes <= 0 {
		return nil, xerrors.New("readBonesWorldPos: no hitboxes present")
	}

	var hitboxIndex = editor.Read2(4, studioHitboxSet+0x8).Int32()
	var start int32 = 0
	it := start

	var Bones = make(map[uintptr]memory.Vector3)
	for it = 0; it < numHitboxes; it++ {
		sz := int32(unsafe.Sizeof(memory.HitBox{}))
		studioHdrData := editor.Read2(sz, 0x44*(it-start)+hitboxIndex+studioHitboxSet).HitBox()

		sz = int32(unsafe.Sizeof(memory.CSMatrix{}))
		theBone := editor.Read2(sz, bonePre+0x30*studioHdrData.Bone).CSMatrix()

		var bonePos = memory.Vector3{theBone[0][3], theBone[1][3], theBone[2][3]}

		Bones[uintptr(it)] = bonePos
	}
	return Bones, nil
}

func cleanupName(s string) string {
	runes := []rune(s)
	newName := make([]rune, len(runes))
	for _, char := range s {
		if unicode.IsPrint(char) {
			newName = append(newName, char)
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
		INFO_OFFSET+playerIndex*ENTRY_SIZE, 0)

	if ferror.ErrIsImportant(err) {
		log.Println("getName:", err)
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