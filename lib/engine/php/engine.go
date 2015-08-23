package php

import (
	// Standard library
	"fmt"
	"io"
	"runtime"

	// Internal packages
	"github.com/deuill/sigil/lib/engine"
)

type Engine struct{}

func (e *Engine) NewContext(w io.Writer) (engine.Context, error) {
	ctx := &Context{
		engine: e,
		writer: w,
	}

	return ctx, nil
}

func init() {
	fmt.Println("Initializing engine...")
	e := &Engine{}

	runtime.SetFinalizer(e, func(e *Engine) {
		fmt.Println("Destroying engine...")
	})

	engine.Register("php", e, []string{".php"})
}
