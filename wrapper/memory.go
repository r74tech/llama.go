package wrapper

/*
#include "../core/include/process.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/Qitmeer/llama.go/config"
)

// LoadFromMemory loads a model from a memory buffer
func LoadFromMemory(modelData []byte, cfg *config.Config) error {
	if len(modelData) == 0 {
		return fmt.Errorf("empty model data")
	}

	// Pin memory to prevent GC from moving it
	var pin runtime.Pinner
	pin.Pin(&modelData[0])
	defer pin.Unpin()

	// Build configuration arguments
	args := buildConfigArgs(cfg)
	cargs := C.CString(args)
	defer C.free(unsafe.Pointer(cargs))

	// Initial prompt
	cprompt := C.CString(cfg.Prompt)
	defer C.free(unsafe.Pointer(cprompt))

	// Call the C function to load from memory
	ret := C.llama_start_from_memory(
		unsafe.Pointer(&modelData[0]),
		C.size_t(len(modelData)),
		cargs,
		C.int(0), // async = false
		cprompt,
	)

	if ret != 0 {
		return fmt.Errorf("failed to load model from memory")
	}

	return nil
}

// LoadFromMmap loads a model from memory-mapped data
func LoadFromMmap(addr uintptr, data []byte, cfg *config.Config) error {
	if len(data) == 0 {
		return fmt.Errorf("empty mmap data")
	}

	// Build configuration arguments
	args := buildConfigArgs(cfg)
	cargs := C.CString(args)
	defer C.free(unsafe.Pointer(cargs))

	// Initial prompt
	cprompt := C.CString(cfg.Prompt)
	defer C.free(unsafe.Pointer(cprompt))

	// Call the C function to load from mmap
	ret := C.llama_start_from_mmap(
		unsafe.Pointer(addr),
		C.size_t(len(data)),
		cargs,
		C.int(0), // async = false
		cprompt,
	)

	if ret != 0 {
		return fmt.Errorf("failed to load model from mmap")
	}

	return nil
}

// buildConfigArgs builds configuration arguments string for C++ side
func buildConfigArgs(cfg *config.Config) string {
	args := ""

	// Model path (for reference, not actually used in memory loading)
	if cfg.Model != "" {
		args += fmt.Sprintf(" --model %s", cfg.Model)
	}

	// Context size
	if cfg.CtxSize > 0 {
		args += fmt.Sprintf(" --ctx-size %d", cfg.CtxSize)
	}

	// GPU layers
	if cfg.NGpuLayers > 0 {
		args += fmt.Sprintf(" --n-gpu-layers %d", cfg.NGpuLayers)
	}

	// Number of predictions
	if cfg.NPredict > 0 {
		args += fmt.Sprintf(" --n-predict %d", cfg.NPredict)
	}

	// Seed
	if cfg.Seed > 0 {
		args += fmt.Sprintf(" --seed %d", cfg.Seed)
	}

	// Batch size
	if cfg.BatchSize > 0 {
		args += fmt.Sprintf(" --batch-size %d", cfg.BatchSize)
	}

	// Ubatch size
	if cfg.UBatchSize > 0 {
		args += fmt.Sprintf(" --ubatch-size %d", cfg.UBatchSize)
	}

	// Pooling type for embeddings
	if cfg.Pooling != "" {
		args += fmt.Sprintf(" --pooling %s", cfg.Pooling)
	}

	// Interactive mode
	if cfg.Interactive {
		args += " --interactive"
	}

	return args
}
