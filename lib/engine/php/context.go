package php

// #include <stdio.h>
// #include <stdlib.h>
// #include "engine.h"
// #include "context.h"
import "C"

import (
	// Standard library
	"fmt"
	"io"
	"unsafe"
)

type Context struct {
	context *C.struct__engine_context
	writer  io.Writer
}

func (c *Context) Bind(name string, value interface{}) error {
	fmt.Printf("Binding variable '%#v\n to name '%s'.", value, name)

	return nil
}

func (c *Context) Run(filename string) error {
	name := C.CString(filename)
	defer C.free(unsafe.Pointer(name))

	_, err := C.context_run(c.context, name)
	if err != nil {
		return fmt.Errorf("Error executing script '%s' in context", filename)
	}

	return nil
}

//export ContextWrite
func ContextWrite(ctxptr unsafe.Pointer, buffer unsafe.Pointer, length C.int) C.int {
	context := (*Context)(ctxptr)

	written, err := context.writer.Write(C.GoBytes(buffer, length))
	if err != nil {
		return C.int(0)
	}

	return C.int(written)
}
