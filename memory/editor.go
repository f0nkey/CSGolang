package memory

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"github.com/AllenDang/w32"
	"log"
	"reflect"
	"unsafe"
)

var (
	procWriteProcessMemory = modkernel32.NewProc("WriteProcessMemory")
	procReadProcessMemory  = modkernel32.NewProc("ReadProcessMemory")
)

type Editor struct {
	procHandle w32.HANDLE
	byteOrder  binary.ByteOrder
	DLLEngine uintptr
	DLLClient uintptr
}

func NewEditor(procName string) (e *Editor) {

	pid, err := FindProcessByName(procName)
	if err != nil {
		log.Fatal(err)
	}

	procHandle, err := w32.OpenProcess(w32.PROCESS_ALL_ACCESS, false, uintptr(pid))
	if err != nil {
		log.Fatal(err)
	}

	dllEngine := uintptr(unsafe.Pointer(GetDLLModuleAddress("engine.dll", pid)))
	dllClient := uintptr(unsafe.Pointer(GetDLLModuleAddress("client_panorama.dll", pid)))

	return &Editor{procHandle, binary.LittleEndian, dllEngine,dllClient}
}




func (e Editor) Read(size int32, offsets ...uintptr) (RawData, error) {
	var err error
	var offset uintptr
	for i, currOffset := range offsets {
		offset += currOffset
		buf := make([]byte, size)
		_, _, err := procReadProcessMemory.Call(
			uintptr(e.procHandle), //handle to dll within proc
			uintptr(offset),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(size),                     //size of datatype
			uintptr(unsafe.Pointer(nil)), //bytesRead
		)

		if i == len(offsets)-1 { // no longer must dereference, reached final value
			return buf, err
		}

		var m int32
		newReader := bytes.NewBuffer(buf)
		if err := binary.Read(newReader, binary.LittleEndian, &m); err != nil {
			log.Println(err)
		}
		offset = uintptr(m)

	}

	return RawData{}, err
}

func (e Editor) Write(addr uintptr, data interface{}) error {
	typeSize := reflect.TypeOf(data).Size()

	//postData, err := GetBytes(data)
	//if err != nil {
	//	log.Println("Editor.Write: can not convert data to bytes")
	//}
	//buf := data //first four are the interface header, useless


	preBuf := &bytes.Buffer{}
	err := binary.Write(preBuf, binary.LittleEndian, data)
	if err != nil {
		return err
	}

	buf := preBuf.Bytes()
	err1, _, status := procWriteProcessMemory.Call(
		uintptr(e.procHandle),
		addr,
		uintptr(unsafe.Pointer(&buf[0])), //lpBuffer
		typeSize,
		uintptr(unsafe.Pointer(nil)),
	)
	if err1 != 1 {
		return status
	}

	return nil
}

type RawData []byte

func (r RawData) Int32() int32 {
	return int32(binary.LittleEndian.Uint32(r))
}

func (r RawData) Uintptr() uintptr {
	return uintptr(int32(binary.LittleEndian.Uint32(r)))

}

func (r RawData) String() string {
	s := string(r[:32])
	return s
}

func (r RawData) HitBox() HitBox {
	var hb HitBox
	newReader := bytes.NewBuffer(r)
	binary.Read(newReader, binary.LittleEndian, &hb)
	return hb
}

func (r RawData) CSMatrix() CSMatrix {
	var vm CSMatrix
	newReader := bytes.NewBuffer(r)
	binary.Read(newReader, binary.LittleEndian, &vm)
	return vm
}

type CSMatrix [4][4]float32

type HitBox struct {
	Bone     int32      /// 0x00
	Group    int32      /// 0x04
	Mins     Vector3    /// 0x08
	Maxs     Vector3    /// 0x14
	HitBoxID int32      /// 0x20
	Pad1     [0xC]byte  /// 0x24
	Radius   float32    /// 0x30
	Pad2     [0x10]byte /// 0x34
}

const (
	Hitbox_Head            = iota
	Hitbox_Neck            = iota
	Hitbox_Pelvis          = iota
	Hitbox_Stomach         = iota
	Hitbox_Chest           = iota
	Hitbox_Lower_Chest     = iota
	Hitbox_Upper_Chest     = iota
	Hitbox_Right_Thigh     = iota
	Hitbox_Left_Thigh      = iota
	Hitbox_Right_Calf      = iota
	Hitbox_Left_Calf       = iota
	Hitbox_Right_Foot      = iota
	Hitbox_Left_Foot       = iota
	Hitbox_Right_Hand      = iota
	Hitbox_Left_Hand       = iota
	Hitbox_Right_Upper_Arm = iota
	Hitbox_Right_Lower_Arm = iota
	Hitbox_Left_Upper_Arm  = iota
	Hitbox_Left_Lower_Arm  = iota
	Hitbox_Last            = iota
)

var hitBoxIDS []string = []string{
	"Hitbox_Head",
	"Hitbox_Neck",
	"Hitbox_Pelvis",
	"Hitbox_Stomach",
	"Hitbox_Chest",
	"Hitbox_Lower_Chest",
	"Hitbox_Upper_Chest",
	"Hitbox_Right_Thigh",
	"Hitbox_Left_Thigh",
	"Hitbox_Right_Calf",
	"Hitbox_Left_Calf",
	"Hitbox_Right_Foot",
	"Hitbox_Left_Foot",
	"Hitbox_Right_Hand",
	"Hitbox_Left_Hand",
	"Hitbox_Right_Upper_Arm",
	"Hitbox_Right_Lower_Arm",
	"Hitbox_Left_Upper_Arm",
	"Hitbox_Left_Lower_Arm",
	"Hitbox_Last",
}

type Vector3 struct {
	X, Y, Z float32
}

type Vector2 struct {
	X, Y float32
}

//TODO: properly make Vector3
//func (r RawData) Vector3() uintptr {
//	return uintptr(int32(binary.LittleEndian.Uint32(r)))
//
//}

func (r RawData) InternalPlayer() (p InternalPlayer) {
	newReader := bytes.NewBuffer(r)
	if err := binary.Read(newReader, binary.LittleEndian, &p); err != nil {
		log.Println("RawData.InternalPlayer:",err)
	}
	return p
}

type InternalPlayer struct {
	//Pad1 [0xF4]byte //0x0000
	Team int32      //0x00F4
	Pad1 [0x8]byte  //0x00F8 oHealth-oTeamNumber-4
	HP   int32      //0x0100
	Pad2 [0x34]byte //0x0104 oVecOrigin-oHealth-4
	X    float32    //0x0138
	Y    float32    //0x013C
	Z    float32    //0x0140
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}
