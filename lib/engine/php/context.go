package php

import (
	// Standard library
	"fmt"
	"io"
)

type Context struct {
	engine *Engine
	writer io.Writer
}

func (c *Context) Bind(name string, value interface{}) error {
	fmt.Printf("Binding variable '%#v\n to name '%s'.", value, name)
	return nil
}

func (c *Context) Run(file string) error {
	fmt.Printf("Running script %s\n", file)
	return nil
}
