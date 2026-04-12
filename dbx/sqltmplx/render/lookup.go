package render

import (
	"reflect"
	"strings"

	"github.com/DaiYuANg/arcgo/collectionx"
	"github.com/samber/lo"
	"github.com/samber/mo"
)

func lookup(params any, name string) mo.Option[any] {
	cur := params
	remaining := name
	for remaining != "" {
		part, rest, found := strings.Cut(remaining, ".")
		remaining = rest

		next := lookupOne(cur, part)
		if next.IsAbsent() {
			return mo.None[any]()
		}
		cur = next.MustGet()
		if !found {
			break
		}
	}
	return mo.Some(cur)
}

func lookupOne(params any, name string) mo.Option[any] {
	if provider, ok := params.(paramLookup); ok {
		value, exists := provider.LookupSQLTemplateParam(name)
		if exists {
			return mo.Some(value)
		}
	}
	v, ok := indirectValue(params)
	if !ok {
		return mo.None[any]()
	}
	if v.Kind() == reflect.Map {
		return lookupMapValue(v, name)
	}
	if v.Kind() == reflect.Struct {
		return lookupStructValue(v, name)
	}
	return mo.None[any]()
}

func lookupStructValue(v reflect.Value, name string) mo.Option[any] {
	meta := cachedStructMetadata(v.Type())
	if field, exists := meta.lookup.Get(strings.ToLower(name)); exists {
		return mo.Some(v.Field(field.index).Interface())
	}
	if value, ok := callZeroArgMethod(v, name); ok {
		return mo.Some(value)
	}
	return mo.None[any]()
}

func callZeroArgMethod(v reflect.Value, name string) (any, bool) {
	if value, ok := callZeroArgMethodOn(v, name); ok {
		return value, true
	}
	if v.Kind() != reflect.Pointer && v.CanAddr() {
		return callZeroArgMethodOn(v.Addr(), name)
	}
	return nil, false
}

func callZeroArgMethodOn(v reflect.Value, name string) (any, bool) {
	for index := range v.Type().NumMethod() {
		method := v.Type().Method(index)
		if !strings.EqualFold(method.Name, name) {
			continue
		}
		value := v.Method(index)
		if value.Type().NumIn() != 0 || value.Type().NumOut() != 1 {
			continue
		}
		return value.Call(nil)[0].Interface(), true
	}
	return nil, false
}

func lookupMapValue(v reflect.Value, name string) mo.Option[any] {
	if v.Type().Key().Kind() != reflect.String {
		return mo.None[any]()
	}
	if value, ok := reflectMapStringValue(v, name); ok {
		return mo.Some(value)
	}
	if value, ok := reflectMapStringValue(v, strings.ToLower(name)); ok {
		return mo.Some(value)
	}
	if value, ok := reflectMapStringValue(v, strings.ToUpper(name)); ok {
		return mo.Some(value)
	}
	return mo.None[any]()
}

func reflectMapStringValue(v reflect.Value, key string) (any, bool) {
	mapKey := reflect.ValueOf(key)
	if keyType := v.Type().Key(); mapKey.Type() != keyType && mapKey.Type().ConvertibleTo(keyType) {
		mapKey = mapKey.Convert(keyType)
	}
	mv := v.MapIndex(mapKey)
	if !mv.IsValid() {
		return nil, false
	}
	return mv.Interface(), true
}

func fieldAliases(f reflect.StructField) collectionx.List[string] {
	aliases := collectionx.NewListWithCapacity[string](3)
	seen := collectionx.NewSetWithCapacity[string](3)
	for _, tagKey := range [...]string{"sqltmpl", "db", "json"} {
		raw := strings.TrimSpace(f.Tag.Get(tagKey))
		if raw == "" || raw == "-" {
			continue
		}
		alias := strings.TrimSpace(strings.Split(raw, ",")[0])
		if alias == "" || alias == "-" {
			continue
		}
		if seen.Contains(alias) {
			continue
		}
		seen.Add(alias)
		aliases.Add(alias)
	}
	return aliases
}

func isEmpty(v any) bool {
	if v == nil {
		return true
	}
	rv, ok := indirectValue(v)
	if !ok {
		return true
	}
	if isLengthValue(rv.Kind()) {
		return rv.Len() == 0
	}
	return false
}

func isLengthValue(kind reflect.Kind) bool {
	return kind == reflect.String || kind == reflect.Array || kind == reflect.Slice || kind == reflect.Map
}

func isBlank(v any) bool {
	if v == nil {
		return true
	}
	rv, ok := indirectValue(v)
	if !ok {
		return true
	}
	return lo.IsEmpty(rv.Interface())
}

func isPresent(v any) bool {
	return !isBlank(v)
}
