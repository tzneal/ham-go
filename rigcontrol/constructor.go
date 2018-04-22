package rigcontrol

import (
	"fmt"
	"reflect"
)

func NewRig(typ string, cfg map[string]interface{}) (Rig, error) {
	switch typ {
	case "FT857D":
		return newft857d(cfg)
	default:
		return nil, fmt.Errorf("unsupported rig type: %s", typ)
	}
	/*
			rig, err := rigcontrol.NewFT857D(rigcontrol.FT857DOptions{
			Port:     "/dev/ttyUSB0",
			BaudRate: 4800,
			DataBits: 8,
			StopBits: 2,
		})
		_ = rig
		_ = err

	*/

	return nil, nil
}

func verifyProperty(cfg map[string]interface{}, name string, t reflect.Type) (bool, error) {
	value, hasProp := cfg[name]
	if !hasProp {
		return false, fmt.Errorf("property %s of type %s must be set", name, t)
	}
	if reflect.TypeOf(value) != t {
		return false, fmt.Errorf("property %s has type %s but must be of type %s", name, reflect.TypeOf(value), t)
	}
	return true, nil
}

func newft857d(cfg map[string]interface{}) (Rig, error) {
	if ok, err := verifyProperty(cfg, "port", reflect.TypeOf("")); !ok {
		return nil, err
	}
	if ok, err := verifyProperty(cfg, "baudrate", reflect.TypeOf(int64(0))); !ok {
		return nil, err
	}
	if ok, err := verifyProperty(cfg, "databits", reflect.TypeOf(int64(0))); !ok {
		return nil, err
	}
	if ok, err := verifyProperty(cfg, "stopbits", reflect.TypeOf(int64(0))); !ok {
		return nil, err
	}

	opts := FT857DOptions{
		Port:     cfg["port"].(string),
		BaudRate: uint(cfg["baudrate"].(int64)),
		DataBits: uint(cfg["databits"].(int64)),
		StopBits: uint(cfg["stopbits"].(int64)),
	}
	return NewFT857D(opts)
}
