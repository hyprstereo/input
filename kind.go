package input

import (
	"fmt"
	"reflect"
	"strconv"
)

type Kind = uint8

const (
	Null Kind = iota
	Bool
	Int
	// Int8
	// Int16
	// Int32
	// Int64
	Uint
	// Uint8
	// Uint16
	// Uint32
	// Uint64
	Float
	// Float32
	// Float64
	Byte
	String
	Map
	Array
	Any
	RGBHex
)

func KindString(typ Kind) (str string) {
	switch typ {
	case RGBHex:
		str = "RGBHex"
	case Null:
		str = "Null"
	case Int:
		str = "Int"
	case Byte:
		str = "Uint32"
	case Bool:
		str = "Bool"
	case Float:
		str = "Float"
	case Map:
		str = "Map"
	case Array:
		str = "Array"
	case Any:
		str = "Any"
	default:
		str = fmt.Sprint(typ)
	}
	return
}

func kindOf(v any) (s Kind) {
	tv := reflect.TypeOf(v)
	switch tv.Kind() {
	case reflect.Array:
		s = Array
	case reflect.Bool:
		s = Bool
	case reflect.Float64:
		s = Float
	case reflect.Float32:
		s = Float
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s = Int
	case reflect.Uint, reflect.Uint16, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		s = Uint
	case reflect.Map:
		s = Map
	case reflect.Interface:
		s = Any
	case reflect.Slice:
		s = Byte
	case reflect.String:
		s = String
	default:
		s = uint8(tv.Kind())
	}
	return
}

func StringToKind(typ string) (str Kind) {
	switch typ {
	case "String":
		str = String
	case "RGBHex":
		str = RGBHex
	case "Null":
		str = Null
	case "Int":
		str = Int
	case "Uint":
		str = Uint
	case "Byte":
		str = Byte
	case "Bool":
		str = Bool
	case "Float":
		str = Float
	case "Float32":
		str = Float

	case "Map":
		str = Map
	case "Array":
		str = Array
	case "Any":
		str = Any
	}
	return
}

func KindFmtSymbol(k Kind) (s string) {
	switch k {
	case Int:
		s = "%d"
	case Float:
		s = "%f"
	case String:
		s = "%s"
	case Bool:
		s = "%t"
	case Any, Map, Array, Byte:
		s = "%v"
	case RGBHex:
		s = "%02x%02x%02x"
	}
	return
}

func KindValue(k Kind) (sy any) {
	var s any
	switch k {
	case Int, Uint:
		s = 0
	case Float:
		s = 0.0
	case String:
		s = ""
	case Bool:
		s = true
	case Array:
		s = []any{}
	case Any, Map:
		s = map[string]any{}
	case RGBHex:
		s = "%02x%02x%02x"
	case Byte:
		s = []byte("")
	}
	sy = reflect.New(reflect.ValueOf(s).Type()).Interface()
	return
}

func pointerValue(v any, k Kind) (val any) {
	va := reflect.ValueOf(v)
	if va.Kind() == reflect.Ptr {
		val = va.Elem().Interface()
	} else {
		val = va.Interface()
	}

	return
}

type Var struct {
	Pos          int
	Value        any
	Name         string
	fmtValue     string
	expectedKind Kind
}

func (v *Var) String() string {
	return fmt.Sprintf("%s (%s): %v", v.Name, KindString(v.expectedKind), v.Value)
}

func (v *Var) tryConvert(value any) (err error) {
	switch v.expectedKind {
	case String:
		v.Value = value
	case Int:
		if va, er := strconv.ParseInt(value.(string), 0, 0); er != nil {
			v.Value = "not_" + KindString(v.expectedKind)
		} else {
			v.Value = va
		}
	case Float:
		if va, er := strconv.ParseFloat(value.(string), 64); er != nil {
			v.Value = "not_" + KindString(v.expectedKind)
		} else {
			v.Value = va
		}
	case Uint:
		if va, er := strconv.ParseUint(value.(string), 0, 0); er != nil {
			v.Value = "not_" + KindString(v.expectedKind)
		} else {
			v.Value = va
		}
	case Bool:
		if va, er := strconv.ParseBool(value.(string)); er != nil {
			v.Value = "not_" + KindString(v.expectedKind)
		} else {
			v.Value = va
		}
	default:
		v.Value = value
	}
	return
}

func (v *Var) Type() (k reflect.Kind) {
	k = reflect.TypeOf(v.Value).Kind()
	return
}
