package httpx

import (
	"github.com/DaiYuANg/arcgo/pkg/option"
	"github.com/danielgtaylor/huma/v2"
)

func buildOperationMutation(operationOptions []OperationOption) func(*huma.Operation) {
	return func(op *huma.Operation) {
		if op == nil {
			return
		}
		option.Apply(op, operationOptions...)
	}
}

func applyOperationMutations(op *huma.Operation, mutators []func(*huma.Operation)) {
	if op == nil || len(mutators) == 0 {
		return
	}
	option.Apply(op, mutators...)
}

func applyWrappers[T any](handler T, wrappers []func(T) T) T {
	wrapped := handler
	for i := len(wrappers) - 1; i >= 0; i-- {
		if wrappers[i] != nil {
			wrapped = wrappers[i](wrapped)
		}
	}
	return wrapped
}
