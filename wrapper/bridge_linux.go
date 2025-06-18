//go:build linux && !cuda

package wrapper

/*
#cgo CFLAGS: -std=c11
#cgo CXXFLAGS: -std=c++17
#cgo CFLAGS: -I${SRCDIR}/../core/include
#cgo CXXFLAGS: -I${SRCDIR}/../core/include
#cgo LDFLAGS: -L${SRCDIR}/../build/lib -lllama_core -lllama -lcommon -lggml -lggml-base -lggml-cpu -lstdc++ -lm
#include <stdlib.h>
#include "core.h"
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/Qitmeer/llama.go/config"
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
	C.free(unsafe.Pointer(ret))
	return content, nil
}

func LlamaInteractive(cfg *config.Config) error {
	if !cfg.Interactive {
		return fmt.Errorf("Not config interactive")
	}
	if len(cfg.Model) <= 0 {
		return fmt.Errorf("No model")
	}
	ip := C.CString(cfg.Prompt)
	defer C.free(unsafe.Pointer(ip))

	cfgArgs := fmt.Sprintf("llama -i --model %s --ctx-size %d --n-gpu-layers %d --n-predict %d --seed %d",
		cfg.Model, cfg.CtxSize, cfg.NGpuLayers, cfg.NPredict, cfg.Seed)
	ca := C.CString(cfgArgs)
	defer C.free(unsafe.Pointer(ca))

	ret := C.llama_start(ca, 0, ip)
	if ret != 0 {
		return fmt.Errorf("Llama start error")
	}
	ret = C.llama_stop()
	if ret != 0 {
		return fmt.Errorf("Llama stop error")
	}
	return nil
}

func LlamaProcess(prompt string) (string, error) {
	if len(prompt) <= 0 {
		return "", fmt.Errorf("No prompt")
	}
	ip := C.CString(prompt)
	defer C.free(unsafe.Pointer(ip))

	ret := C.llama_gen(ip)
	if ret == nil {
		return "", fmt.Errorf("Llama run error")
	}
	content := C.GoString(ret)
	C.free(unsafe.Pointer(ret))
	return content, nil
}

func LlamaStart(cfg *config.Config) error {
	if len(cfg.Model) <= 0 {
		return fmt.Errorf("No model")
	}
	cfgArgs := fmt.Sprintf("llama -i --model %s --ctx-size %d --n-gpu-layers %d --n-predict %d --seed %d",
		cfg.Model, cfg.CtxSize, cfg.NGpuLayers, cfg.NPredict, cfg.Seed)
	ca := C.CString(cfgArgs)
	defer C.free(unsafe.Pointer(ca))

	ip := C.CString(cfg.Prompt)
	defer C.free(unsafe.Pointer(ip))

	ret := C.llama_start(ca, 1, ip)
	if ret != 0 {
		return fmt.Errorf("Llama start error")
	}
	return nil
}

func LlamaStop() error {
	ret := C.llama_stop()
	if ret != 0 {
		return fmt.Errorf("Llama stop error")
	}
	return nil
}

func LlamaEmbedding(cfg *config.Config) (string, error) {
	if len(cfg.Model) <= 0 {
		return "", fmt.Errorf("No model")
	}
	if len(cfg.Prompt) <= 0 {
		return "", fmt.Errorf("No prompt")
	}
	ip := C.CString(cfg.Prompt)
	defer C.free(unsafe.Pointer(ip))

	cfgArgs := fmt.Sprintf("llama --model %s --ctx-size %d --n-gpu-layers %d --n-predict %d --seed %d --embd-normalize %d --batch-size %d --ubatch-size %d",
		cfg.Model, cfg.CtxSize, cfg.NGpuLayers, cfg.NPredict, cfg.Seed, cfg.EmbdNormalize, cfg.BatchSize, cfg.UBatchSize)
	if len(cfg.Pooling) > 0 {
		cfgArgs = fmt.Sprintf("%s --pooling %s", cfgArgs, cfg.Pooling)
	}
	if len(cfg.EmbdOutputFormat) > 0 {
		cfgArgs = fmt.Sprintf("%s --embd-output-format %s", cfgArgs, cfg.EmbdOutputFormat)
	}
	if len(cfg.EmbdSeparator) > 0 {
		cfgArgs = fmt.Sprintf("%s --embd-separator %s", cfgArgs, cfg.EmbdOutputFormat)
	}
	ca := C.CString(cfgArgs)
	defer C.free(unsafe.Pointer(ca))

	ret := C.llama_embedding(ca, ip)
	if ret == nil {
		return "", fmt.Errorf("llama_embedding run error")
	}
	content := C.GoString(ret)
	C.free(unsafe.Pointer(ret))
	return content, nil
}
