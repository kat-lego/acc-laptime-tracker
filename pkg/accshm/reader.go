package accshm

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	procOpenFileMapping = kernel32.NewProc("OpenFileMappingW")
	procMapViewOfFile   = kernel32.NewProc("MapViewOfFile")
	procUnmapViewOfFile = kernel32.NewProc("UnmapViewOfFile")
	procCloseHandle     = kernel32.NewProc("CloseHandle")
)

const (
	FILE_MAP_READ = 0x0004
)

func ReadSharedMemoryStruct[T any](name string) (*T, error) {
	namePtr, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return nil, fmt.Errorf("failed to encode name: %w\n", err)
	}

	hMap, _, callErr := procOpenFileMapping.Call(
		FILE_MAP_READ,
		0,
		uintptr(unsafe.Pointer(namePtr)),
	)
	if hMap == 0 {
		return nil, fmt.Errorf("OpenFileMappingW failed: %w\n", callErr)
	}
	defer procCloseHandle.Call(hMap)

	structSize := unsafe.Sizeof(new(T))
	addr, _, callErr := procMapViewOfFile.Call(
		hMap,
		FILE_MAP_READ,
		0,
		0,
		uintptr(structSize),
	)
	if addr == 0 {
		return nil, fmt.Errorf("MapViewOfFile failed: %w\n", callErr)
	}
	defer procUnmapViewOfFile.Call(addr)

	value := *(*T)(unsafe.Pointer(addr))
	return &value, nil
}
