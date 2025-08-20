//go:build linux
// +build linux

package wrapper

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// MmapModel maps a model file into memory using mmap on Linux
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

	// Get file descriptor
	fd := int(file.Fd())

	// Perform mmap with proper alignment
	pageSize := syscall.Getpagesize()
	offset := int64(0)

	// Ensure offset is page-aligned
	pageAlignedOffset := (offset / int64(pageSize)) * int64(pageSize)
	adjustment := int(offset - pageAlignedOffset)
	mapSize := size + adjustment

	// Map the file into memory
	// Use PROT_READ for read-only access, MAP_PRIVATE for copy-on-write semantics
	mappedData, err := syscall.Mmap(fd, pageAlignedOffset, mapSize,
		syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		return 0, nil, fmt.Errorf("mmap failed: %w", err)
	}

	// Adjust for any offset within the page
	modelData := mappedData[adjustment:]

	// Get the address of the mapped region
	dataAddr := uintptr(unsafe.Pointer(&modelData[0]))

	// Verify GGUF magic number
	if len(modelData) < 4 {
		syscall.Munmap(mappedData)
		return 0, nil, fmt.Errorf("model file too small to be valid GGUF")
	}

	// Check for GGUF magic: 'G' 'G' 'U' 'F' (0x46554747 in little-endian)
	if modelData[0] != 'G' || modelData[1] != 'G' || modelData[2] != 'U' || modelData[3] != 'F' {
		syscall.Munmap(mappedData)
		return 0, nil, fmt.Errorf("invalid GGUF magic number: got %x %x %x %x",
			modelData[0], modelData[1], modelData[2], modelData[3])
	}

	return dataAddr, mappedData, nil
}

// UnmapModel unmaps a previously memory-mapped model
func UnmapModel(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	err := syscall.Munmap(data)
	if err != nil {
		return fmt.Errorf("munmap failed: %w", err)
	}

	return nil
}

// MmapModelAtOffset maps a model file at a specific offset
func MmapModelAtOffset(fd int, offset int64, size int) (addr uintptr, data []byte, err error) {
	if size <= 0 {
		return 0, nil, fmt.Errorf("invalid size: %d", size)
	}

	// Get page size for alignment
	pageSize := int64(syscall.Getpagesize())

	// Calculate page-aligned offset
	pageAlignedOffset := (offset / pageSize) * pageSize
	adjustment := int(offset - pageAlignedOffset)
	mapSize := size + adjustment

	// Perform the mmap
	mappedData, err := syscall.Mmap(fd, pageAlignedOffset, mapSize,
		syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		return 0, nil, fmt.Errorf("mmap at offset failed: %w", err)
	}

	// Adjust to the actual model data start
	modelData := mappedData[adjustment:]
	dataAddr := uintptr(unsafe.Pointer(&modelData[0]))

	return dataAddr, mappedData, nil
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
