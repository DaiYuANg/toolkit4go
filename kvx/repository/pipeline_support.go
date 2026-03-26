package repository

import (
	"errors"

	"github.com/DaiYuANg/arcgo/kvx"
)

type pipelineProvider interface {
	Pipeline() kvx.Pipeline
}

var ErrExpiration = errors.New("expiration <= 0")
