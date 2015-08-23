package engine

import (
	// Standard library
	"fmt"
	"io"
	"path/filepath"
	"reflect"
)

type Engine interface {
	NewContext(w io.Writer) (Context, error)
}

type Context interface {
	Bind(name string, value interface{}) error
	Run(file string) error
}

type Value interface {
	Int() int64
	Float() float64
	Bool() bool
	String() string
	Kind() reflect.Kind
}

// A map of extensions to Engines that implement them.
var engines map[string]Engine

func Register(name string, rcvr Engine, types []string) error {
	for _, t := range types {
		if _, exists := engines[t]; exists {
			return fmt.Errorf("Extension '%s' already exists, refusing to overwrite", t)
		}

		engines[t] = rcvr
	}

	return nil
}

func Handle(w io.Writer, file string) error {
	ext := filepath.Ext(file)

	engine, exists := engines[ext]
	if !exists {
		return fmt.Errorf("Extension '%s' not registered", ext)
	}

	ctx, err := engine.NewContext(w)
	if err != nil {
		return err
	}

	if err = ctx.Run(file); err != nil {
		return err
	}

	return nil
}

func init() {
	engines = make(map[string]Engine)
}
