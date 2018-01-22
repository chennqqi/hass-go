package dynamic

import (
	"fmt"
	"reflect"
	"time"
)

type Dynamic struct {
	Item interface{}
}

func (dyn Dynamic) Print() {
	fmt.Printf("Type=%v\n", reflect.TypeOf(dyn.Item).String())
}
func (dyn Dynamic) Type() string {
	return reflect.TypeOf(dyn.Item).String()
}
func (dyn Dynamic) IsNil() bool {
	return dyn.Item == nil
}
func (dyn Dynamic) IsArray() bool {
	kind := reflect.TypeOf(dyn.Item).Kind()
	return kind == reflect.Array || kind == reflect.Slice
}
func (dyn Dynamic) IsMap() bool {
	return reflect.TypeOf(dyn.Item).Kind() == reflect.Map
}
func (dyn Dynamic) IsString() bool {
	return reflect.TypeOf(dyn.Item).Kind() == reflect.String
}
func (dyn Dynamic) AsString() string {
	if dyn.IsString() {
		return dyn.Item.(string)
	}
	return ""
}

func (dyn Dynamic) IsDuration() bool {
	if reflect.TypeOf(dyn.Item).Kind() == reflect.String {
		_, e := time.ParseDuration(dyn.AsString())
		return e == nil
	}
	return false
}
func (dyn Dynamic) AsDuration() time.Duration {
	if dyn.IsString() {
		d, e := time.ParseDuration(dyn.AsString())
		if e == nil {
			return d
		}
	}
	return time.Duration(0)
}
func (dyn Dynamic) IsBool() bool {
	return reflect.TypeOf(dyn.Item).Kind() == reflect.Bool
}
func (dyn Dynamic) AsBool() bool {
	return dyn.Item.(bool)
}
func (dyn Dynamic) IsInt32() bool {
	return reflect.TypeOf(dyn.Item).Kind() == reflect.Int
}
func (dyn Dynamic) AsInt32() int32 {
	return dyn.Item.(int32)
}
func (dyn Dynamic) IsInt64() bool {
	return reflect.TypeOf(dyn.Item).Kind() == reflect.Int64
}
func (dyn Dynamic) AsInt64() int64 {
	return dyn.Item.(int64)
}
func (dyn Dynamic) IsFloat32() bool {
	return reflect.TypeOf(dyn.Item).Kind() == reflect.Float32
}
func (dyn Dynamic) AsFloat32() float32 {
	if dyn.IsInt32() {
		return float32(dyn.AsInt32())
	} else if dyn.IsInt64() {
		return float32(dyn.AsInt64())
	} else if dyn.IsFloat64() {
		return float32(dyn.AsFloat64())
	}
	return dyn.Item.(float32)
}
func (dyn Dynamic) IsFloat64() bool {
	return reflect.TypeOf(dyn.Item).Kind() == reflect.Float64
}
func (dyn Dynamic) AsFloat64() float64 {
	if dyn.IsInt32() {
		return float64(dyn.AsInt32())
	} else if dyn.IsInt64() {
		return float64(dyn.AsInt64())
	} else if dyn.IsFloat32() {
		return float64(dyn.AsFloat32())
	}
	return dyn.Item.(float64)
}

func (dyn Dynamic) ArrayLen() int {
	a := dyn.Item.([]interface{})
	return len(a)
}
func (dyn Dynamic) ArrayAt(index int) Dynamic {
	if !dyn.IsNil() {
		a := dyn.Item.([]interface{})
		if index < len(a) {
			return Dynamic{a[index]}
		}
	}
	return Dynamic{nil}
}
func (dyn Dynamic) ArrayIter() []Dynamic {
	if !dyn.IsNil() {
		a := dyn.Item.([]interface{})
		b := []Dynamic{}
		for _, i := range a {
			b = append(b, Dynamic{i})
		}
		return b
	}
	return []Dynamic{}
}
func (dyn Dynamic) MapIter() map[string]Dynamic {
	if !dyn.IsNil() {
		m := dyn.Item.(map[string]interface{})
		b := map[string]Dynamic{}
		for k, v := range m {
			b[k] = Dynamic{v}
		}
		return b
	}
	return map[string]Dynamic{}
}
func (dyn Dynamic) Get(key string) Dynamic {
	if !dyn.IsNil() {
		m := dyn.Item.(map[string]interface{})
		item, exists := m[key]
		if exists {
			return Dynamic{item}
		}
	}
	return Dynamic{nil}
}
