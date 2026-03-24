package repository

import (
	"errors"
	"strconv"
	"time"

	"github.com/DaiYuANg/arcgo/kvx"
)

type pipelineProvider interface {
	Pipeline() kvx.Pipeline
}

var ErrExpiration = errors.New("expiration <= 0")

func enqueueExpire(pipe kvx.Pipeline, key string, expiration time.Duration) error {
	if expiration <= 0 {
		return ErrExpiration
	}

	err := pipe.Enqueue("PEXPIRE", []byte(key), []byte(strconv.FormatInt(expiration.Milliseconds(), 10)))
	if err != nil {
		return err
	}
	return nil
}
