package wrapper

/*
#cgo CXXFLAGS: -std=c++17
#cgo CXXFLAGS: -I${SRCDIR}/../llama.cpp/include
#cgo LDFLAGS: -L${SRCDIR}/build/lib -lllama -lstdc++
#include "bridge.h"
*/
import "C"
import "fmt"

func LlamaApp() error {
	ret := C.llama_app()
	if ret != 0 {
		return fmt.Errorf("Llama exit error")
	}
	return nil
}
