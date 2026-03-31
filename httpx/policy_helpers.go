package httpx

import (
	"github.com/danielgtaylor/huma/v2"
)

func buildOperationMutation(operationOptions []OperationOption) func(*huma.Operation) {
	return func(op *huma.Operation) {
		if op == nil {
			return
		}
		applyOptions(op, operationOptions...)
	}
}

func applyOperationMutations(op *huma.Operation, mutators []func(*huma.Operation)) {
	if op == nil || len(mutators) == 0 {
		return
	}
	applyOptions(op, mutators...)
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
