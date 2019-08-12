package memory

import (
	"fmt"
	"github.com/AllenDang/w32"
	"syscall"
	"unsafe"
)

var (
	modkernel32        = syscall.NewLazyDLL("kernel32.dll")
	user32             = syscall.MustLoadDLL("user32.dll")
	procEnumWindows    = user32.MustFindProc("EnumWindows")
	procGetWindowTextW = user32.MustFindProc("GetWindowTextW")
)

func GetProcessName(id uint32) string {
	snapshot := w32.CreateToolhelp32Snapshot(w32.TH32CS_SNAPMODULE, id)
	if snapshot == w32.ERROR_INVALID_HANDLE {
		return "<UNKNOWN>"
	}
	defer w32.CloseHandle(snapshot)

	var me w32.MODULEENTRY32
	me.Size = uint32(unsafe.Sizeof(me))
	if w32.Module32First(snapshot, &me) {
		//		fmt.Println("this is the damn base address:", &me.ModBaseAddr)

		return w32.UTF16PtrToString(&me.SzModule[0])
	}

	return "<UNKNOWN>"
}
func ListProcesses() []uint32 {
	sz := uint32(5000)
	procs := make([]uint32, sz)
	var bytesReturned uint32
	if w32.EnumProcesses(procs, sz, &bytesReturned) {
		return procs[:int(bytesReturned)/4]
	}
	return []uint32{}
}
func FindProcessByName(name string) (uint32, error) {
	for _, pid := range ListProcesses() {
		if GetProcessName(pid) == name {
			return pid, nil
		}
	}
	return 0, fmt.Errorf("unknown process")
}

func GetDLLModuleAddress(nameDLL string, namePID uint32) *uint8 {

	var h w32.HANDLE //HANDLE type
	//var dll uint32 //DWORD type
	var me32 w32.MODULEENTRY32 //stores retrieve info
	h = w32.CreateToolhelp32Snapshot(w32.TH32CS_SNAPMODULE32+w32.TH32CS_SNAPMODULE, namePID)
	me32.Size = uint32(unsafe.Sizeof(me32))
	if h == w32.ERROR_INVALID_HANDLE {
		fmt.Println("Invalid Handle given.")
	}
	w32.Module32First(h, &me32)
	var killLoop bool = false
	for killLoop != true {
		w32.Module32Next(h, &me32)

		if w32.UTF16PtrToString(&me32.SzModule[0]) == nameDLL {
			return me32.ModBaseAddr
			killLoop = true
		}
	}
	fmt.Println("modBaseAddr:", me32.ModBaseAddr)
	return me32.ModBaseAddr
}

func GetDLLModuleAddressSize(nameDLL string, namePID uint32) uint32 {

	var h w32.HANDLE //HANDLE type
	//var dll uint32 //DWORD type
	var me32 w32.MODULEENTRY32 //stores retrieve info
	h = w32.CreateToolhelp32Snapshot(w32.TH32CS_SNAPMODULE32+w32.TH32CS_SNAPMODULE, namePID)
	me32.Size = uint32(unsafe.Sizeof(me32))
	if h == w32.ERROR_INVALID_HANDLE {
		fmt.Println("Invalid Handle given.")
	}
	w32.Module32First(h, &me32)
	var killLoop bool = false
	for killLoop != true {
		w32.Module32Next(h, &me32)

		if w32.UTF16PtrToString(&me32.SzModule[0]) == nameDLL {
			return me32.ModBaseSize
			killLoop = true
		}
	}
	fmt.Println("modBaseAddr:", me32.ModBaseAddr)
	return me32.Size
}

func GetWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FindWindow(title string) (syscall.Handle, error) {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := GetWindowText(h, &b[0], int32(len(b)))
		if err != nil {
			// ignore the error
			return 1 // continue enumeration
		}
		if syscall.UTF16ToString(b) == title {
			// note the window
			hwnd = h
			return 0 // stop enumeration
		}
		return 1 // continue enumeration
	})
	EnumWindows(cb, 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("No window with title '%s' found", title)
	}
	return hwnd, nil
}

func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetLayeredWindowAttributes(hwnd w32.HWND, cr w32.COLORREF, alpha byte, flags uint32) bool {
	moduser32 := syscall.NewLazyDLL("user32.dll")
	procSetLayeredWindowAttributes := moduser32.NewProc("SetLayeredWindowAttributes")
	r0, _,_ := syscall.Syscall6(procSetLayeredWindowAttributes.Addr(), 4, uintptr(hwnd), uintptr(cr), uintptr(alpha), uintptr(flags), 0, 0)

	return r0 != 0
}