package render

import (
	"reflect"
	"strings"
)

func lookup(params any, name string) (any, bool) {
	v := reflect.ValueOf(params)
	for v.IsValid() && v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, false
		}
		v = v.Elem()
	}
	if !v.IsValid() {
		return nil, false
	}
	if v.Kind() == reflect.Map {
		mv := v.MapIndex(reflect.ValueOf(name))
		if mv.IsValid() {
			return mv.Interface(), true
		}
		return nil, false
	}
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}
			if f.Name == name || strings.EqualFold(f.Name, name) {
				return v.Field(i).Interface(), true
			}
		}
	}
	return nil, false
}

func isEmpty(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	for rv.IsValid() && rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return true
		}
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return true
	}
	switch rv.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return rv.Len() == 0
	}
	return false
}
