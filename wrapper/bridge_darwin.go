//go:build darwin
// +build darwin

package wrapper

/*
#cgo CFLAGS: -std=c11
#cgo CXXFLAGS: -std=c++17
#cgo CFLAGS: -I${SRCDIR}/../core
#cgo CXXFLAGS: -I${SRCDIR}/../core
#cgo LDFLAGS: -framework Foundation -framework Metal -framework MetalKit -framework Accelerate -lstdc++
#cgo LDFLAGS: -L${SRCDIR}/../build/lib -lllama_core -lllama -lggml -lggml-base -lggml-cpu -lggml-blas -lggml-metal
#include <stdlib.h>
#include "core.h"
*/
import "C"

import (
	"fmt"
	"github.com/Qitmeer/llama.go/config"
	"unsafe"
)

func LlamaApp(cfg *config.Config) error {
	mp := C.CString(cfg.Model)
	defer C.free(unsafe.Pointer(mp))

	ip := C.CString(cfg.Prompt)
	defer C.free(unsafe.Pointer(ip))

	ret := C.llama_app(mp, ip, C.int(cfg.NGpuLayers), C.int(cfg.NPredict))
	if ret != 0 {
		return fmt.Errorf("Llama exit error")
	}
	return nil
}
