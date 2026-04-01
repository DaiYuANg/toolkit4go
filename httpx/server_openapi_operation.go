package httpx

import (
	"github.com/DaiYuANg/arcgo/collectionx/list"
	"github.com/danielgtaylor/huma/v2"
)

func forEachOperation(doc *huma.OpenAPI, fn func(*huma.Operation)) {
	if doc == nil || fn == nil {
		return
	}
	for _, pathItem := range doc.Paths {
		if pathItem == nil {
			return
		}
		list.NewList(
			pathItem.Get, pathItem.Put, pathItem.Post, pathItem.Delete,
			pathItem.Options, pathItem.Head, pathItem.Patch, pathItem.Trace,
		).Range(func(_ int, op *huma.Operation) bool {
			if op != nil {
				fn(op)
			}
			return true
		})
	}
}

func appendOperationParameter(op *huma.Operation, param *huma.Param) {
	if op == nil || param == nil {
		return
	}
	for _, existing := range op.Parameters {
		if existing != nil && existing.Name == param.Name && existing.In == param.In {
			return
		}
	}
	op.Parameters = append(op.Parameters, cloneParam(param))
}
