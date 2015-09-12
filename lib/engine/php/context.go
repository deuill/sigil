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

	go c.sync()

	_, err := C.context_run(c.context, name)
	if err != nil {
		return fmt.Errorf("Error executing script '%s' in context", filename)
	}

	return nil
}

func (c *Context) sync() {
	var buf *C.char
	var num C.size_t
	var err error

	for {
		if num, err = C.context_sync(c.context, &buf); err != nil {
			break
		} else if num == 0 {
			continue
		}

		fmt.Fprint(c.writer, C.GoString(buf))
		C.free(unsafe.Pointer(buf))
	}
}
