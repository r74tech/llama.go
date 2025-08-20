//go:build windows
// +build windows

package wrapper

import (
	"fmt"
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

var (
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	procCreateFileMapping = kernel32.NewProc("CreateFileMappingW")
	procMapViewOfFile     = kernel32.NewProc("MapViewOfFile")
	procUnmapViewOfFile   = kernel32.NewProc("UnmapViewOfFile")
)

const (
	FILE_MAP_READ  = 0x0004
	PAGE_READONLY  = 0x02
	INVALID_HANDLE = ^uintptr(0)
)

// MmapModel maps a model file into memory using Windows memory mapping
func MmapModel(path string) (addr uintptr, data []byte, err error) {
	// Open the model file
	file, err := os.Open(path)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to open model file: %w", err)
	}
	defer func() {
		if err != nil {
			file.Close()
		}
	}()

	// Get file statistics to determine size
	stat, err := file.Stat()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to stat model file: %w", err)
	}

	size := int(stat.Size())
	if size == 0 {
		return 0, nil, fmt.Errorf("model file is empty")
	}

	// Get Windows handle from file
	handle := syscall.Handle(file.Fd())

	// Create file mapping
	high := uint32(size >> 32)
	low := uint32(size & 0xffffffff)

	mapHandle, _, err := procCreateFileMapping.Call(
		uintptr(handle),
		0, // No security attributes
		PAGE_READONLY,
		uintptr(high),
		uintptr(low),
		0, // No name
	)

	if mapHandle == 0 || mapHandle == INVALID_HANDLE {
		return 0, nil, fmt.Errorf("CreateFileMapping failed: %w", err)
	}
	defer syscall.CloseHandle(syscall.Handle(mapHandle))

	// Map view of file
	mappedAddr, _, err := procMapViewOfFile.Call(
		mapHandle,
		FILE_MAP_READ,
		0, // Offset high
		0, // Offset low
		uintptr(size),
	)

	if mappedAddr == 0 {
		return 0, nil, fmt.Errorf("MapViewOfFile failed: %w", err)
	}

	// Create a byte slice from the mapped memory
	var mappedData []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&mappedData))
	header.Data = mappedAddr
	header.Len = size
	header.Cap = size

	// Verify GGUF magic number
	if len(mappedData) < 4 {
		procUnmapViewOfFile.Call(mappedAddr)
		return 0, nil, fmt.Errorf("model file too small to be valid GGUF")
	}

	// Check for GGUF magic: 'G' 'G' 'U' 'F'
	if mappedData[0] != 'G' || mappedData[1] != 'G' || mappedData[2] != 'U' || mappedData[3] != 'F' {
		procUnmapViewOfFile.Call(mappedAddr)
		return 0, nil, fmt.Errorf("invalid GGUF magic number: got %x %x %x %x",
			mappedData[0], mappedData[1], mappedData[2], mappedData[3])
	}

	return mappedAddr, mappedData, nil
}

// UnmapModel unmaps a previously memory-mapped model
func UnmapModel(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	addr := uintptr(unsafe.Pointer(&data[0]))
	ret, _, err := procUnmapViewOfFile.Call(addr)

	if ret == 0 {
		return fmt.Errorf("UnmapViewOfFile failed: %w", err)
	}

	return nil
}

// MmapModelAtOffset maps a model file at a specific offset
func MmapModelAtOffset(fd int, offset int64, size int) (addr uintptr, data []byte, err error) {
	if size <= 0 {
		return 0, nil, fmt.Errorf("invalid size: %d", size)
	}

	// Get Windows handle
	handle := syscall.Handle(fd)

	// Calculate file mapping parameters
	totalSize := offset + int64(size)
	high := uint32(totalSize >> 32)
	low := uint32(totalSize & 0xffffffff)

	// Create file mapping for the entire needed range
	mapHandle, _, err := procCreateFileMapping.Call(
		uintptr(handle),
		0, // No security attributes
		PAGE_READONLY,
		uintptr(high),
		uintptr(low),
		0, // No name
	)

	if mapHandle == 0 || mapHandle == INVALID_HANDLE {
		return 0, nil, fmt.Errorf("CreateFileMapping failed: %w", err)
	}
	defer syscall.CloseHandle(syscall.Handle(mapHandle))

	// Calculate offset for MapViewOfFile
	offsetHigh := uint32(offset >> 32)
	offsetLow := uint32(offset & 0xffffffff)

	// Map view of file at the specified offset
	mappedAddr, _, err := procMapViewOfFile.Call(
		mapHandle,
		FILE_MAP_READ,
		uintptr(offsetHigh),
		uintptr(offsetLow),
		uintptr(size),
	)

	if mappedAddr == 0 {
		return 0, nil, fmt.Errorf("MapViewOfFile at offset failed: %w", err)
	}

	// Create a byte slice from the mapped memory
	var mappedData []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&mappedData))
	header.Data = mappedAddr
	header.Len = size
	header.Cap = size

	return mappedAddr, mappedData, nil
}

// LoadSelfContainedModel loads a model that's embedded in the current executable
func LoadSelfContainedModel(offset int64, size int) (addr uintptr, data []byte, err error) {
	// Open the current executable
	executable, err := os.Executable()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	file, err := os.Open(executable)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to open executable: %w", err)
	}
	defer file.Close()

	// Map the model portion of the executable
	return MmapModelAtOffset(int(file.Fd()), offset, size)
}
