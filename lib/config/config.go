package config

import (
	// Standard library
	"fmt"
	"reflect"
	"strings"

	// External packages
	ini "github.com/rakyll/goini"
)

type Config struct {
	data ini.Dict
}

func (c *Config) Unpack(rcvr interface{}) error {
	value := reflect.Indirect(reflect.ValueOf(rcvr))

	// Do not attempt to unpack into anything other than a structure.
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("Unsupported receiver type: %s", value.Kind().String())
	}

	return c.unpack("", value)
}

func (c *Config) unpack(section string, value reflect.Value) error {
	// Normalize section name.
	section = strings.ToLower(section)

	// Set each field in struct according to its underlying type.
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		name := strings.ToLower(field.Name)

		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v, exists := c.data.GetInt(section, name); exists {
				value.Field(i).SetInt(int64(v))
			}
		case reflect.Bool:
			if v, exists := c.data.GetBool(section, name); exists {
				value.Field(i).SetBool(v)
			}
		case reflect.String:
			if v, exists := c.data.GetString(section, name); exists {
				value.Field(i).SetString(v)
			}
		case reflect.Struct:
			c.unpack(name, value.Field(i))
		}
	}

	return nil
}

func Load(filename string) (*Config, error) {
	data, err := ini.Load(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{data: data}

	return config, nil
}
