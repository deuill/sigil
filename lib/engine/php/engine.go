package php

// #cgo CFLAGS: -I/usr/include/php -I/usr/include/php/main -I/usr/include/php/TSRM
// #cgo CFLAGS: -I/usr/include/php/Zend
// #cgo LDFLAGS: -lphp5
//
// #include <stdio.h>
// #include "engine.h"
// #include "context.h"
import "C"

import (
	// Standard library
	"fmt"
	"io"
	"runtime"
	"unsafe"

	// Internal packages
	"github.com/deuill/sigil/lib/engine"
)

type Engine struct {
	engine *C.struct__php_engine
}

func (e *Engine) NewContext(w io.Writer) (engine.Context, error) {
	ctx := &Context{writer: w}

	ptr, err := C.context_new(e.engine, unsafe.Pointer(ctx))
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize context for PHP engine")
	}

	ctx.context = ptr

	runtime.SetFinalizer(ctx, func(ctx *Context) {
		C.context_destroy(ctx.context)
	})

	return ctx, nil
}

func init() {
	ptr, err := C.engine_init()
	if err != nil {
		panic("PHP engine failed to initialize")
	}

	e := &Engine{
		engine: ptr,
	}

	runtime.SetFinalizer(e, func(e *Engine) {
		C.engine_shutdown(e.engine)
	})

	engine.Register("php", e, []string{".php"})
}
