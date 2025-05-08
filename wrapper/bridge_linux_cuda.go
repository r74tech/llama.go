//go:build linux && cuda

package wrapper

/*
#cgo CFLAGS: -std=c11
#cgo CXXFLAGS: -std=c++17
#cgo CFLAGS: -I${SRCDIR}/../core
#cgo CXXFLAGS: -I${SRCDIR}/../core
#cgo LDFLAGS: -L${SRCDIR}/../build/lib -lllama_core -lllama -lcommon -lggml -lggml-base -lggml-cpu -lggml-cuda -lstdc++ -lm
#cgo LDFLAGS: -L/usr/local/cuda/lib64 -lcudart -lcublas -L/usr/local/cuda/lib64/stubs -lcuda
#include <stdlib.h>
#include "core.h"
*/
import "C"

import (
	"fmt"
	"github.com/Qitmeer/llama.go/config"
	"unsafe"
)

func LlamaGenerate(cfg *config.Config) (string, error) {
	mp := C.CString(cfg.Model)
	defer C.free(unsafe.Pointer(mp))

	ip := C.CString(cfg.Prompt)
	defer C.free(unsafe.Pointer(ip))

	ret := C.llama_generate(mp, ip, C.int(cfg.NGpuLayers), C.int(cfg.NPredict))
	if ret == nil {
		return "", fmt.Errorf("Llama run error")
	}
	content := C.GoString(ret)
	return content, nil
}

func LlamaInteractive(cfg *config.Config) error {
	mp := C.CString(cfg.Model)
	defer C.free(unsafe.Pointer(mp))

	ret := C.llama_interactive(mp, C.int(cfg.NGpuLayers), C.int(cfg.CtxSize))
	if ret != 0 {
		return fmt.Errorf("Llama exit error")
	}
	return nil
}

func LlamaProcess(cfg *config.Config) (string, error) {
	if cfg.Interactive {
		return "", fmt.Errorf("Not support")
	}
	if len(cfg.Model) <= 0 {
		return "", fmt.Errorf("No model")
	}
	if len(cfg.Prompt) <= 0 {
		return "", fmt.Errorf("No prompt")
	}
	cfgArgs := fmt.Sprintf("llama -no-cnv --model %s --prompt %s --ctx-size %d --n-gpu-layers %d --n-predict %d --seed %d",
		cfg.Model, cfg.Prompt, cfg.CtxSize, cfg.NGpuLayers, cfg.NPredict, cfg.Seed)

	ca := C.CString(cfgArgs)
	defer C.free(unsafe.Pointer(ca))

	ret := C.llama_process(ca)
	if ret == nil {
		return "", fmt.Errorf("Llama run error")
	}
	content := C.GoString(ret)
	C.free(unsafe.Pointer(ret))
	return content, nil
}
