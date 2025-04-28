package wrapper

/*
#cgo CFLAGS: -std=c11
#cgo CXXFLAGS: -std=c++17
#cgo CFLAGS: -I${SRCDIR}/../core
#cgo CXXFLAGS: -I${SRCDIR}/../core
#cgo LDFLAGS: -L${SRCDIR}/../build/lib -lllama_core -lstdc++
#include "core.h"
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
